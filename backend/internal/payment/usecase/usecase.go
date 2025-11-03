package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	customerInterfaces "backend/internal/customer/interfaces"
	models "backend/internal/db_model"
	"backend/internal/payment/dto"
	"backend/internal/payment/interfaces"
	restaurantInterfaces "backend/internal/restaurant/interfaces"
)

type PaymentUsecase struct {
	paymentRepository    interfaces.PaymentRepository
	customerRepository   customerInterfaces.CustomerRepository
	restaurantRepository restaurantInterfaces.RestaurantRepository
}

func NewPaymentUsecase(paymentRepository interfaces.PaymentRepository,
	customerRepository customerInterfaces.CustomerRepository,
	restaurantRepository restaurantInterfaces.RestaurantRepository,
) interfaces.PaymentUsecase {
	return &PaymentUsecase{
		paymentRepository:    paymentRepository,
		customerRepository:   customerRepository,
		restaurantRepository: restaurantRepository,
	}
}

func (u *PaymentUsecase) GetTopupPaymentMethods(userID uuid.UUID) ([]dto.PaymentMethodDetail, error) {

	var paymentMethods []dto.PaymentMethodDetail

	methods, err := u.paymentRepository.GetPaymentMethodsByType("topup")
	if err != nil {
		return nil, err
	}

	for _, method := range methods {
		paymentMethods = append(paymentMethods, dto.PaymentMethodDetail{
			PaymentMethodID: method.ID,
			Name:            method.Name,
			// ImageURL:        method.ImageURL,
		})
	}

	return paymentMethods, nil
}

func (u *PaymentUsecase) TopupToWallet(userID uuid.UUID, request *dto.TopupRequest) error {
	paymentMethod, err := u.paymentRepository.GetPaymentMethodByID(request.PaymentMethodID)
	if err != nil {
		return err
	}
	if paymentMethod.Type != "topup" && paymentMethod.Type != "all" && paymentMethod.Type != "both" {
		return fmt.Errorf("Invalid payment method for top-up")
	}

	customer, err := u.customerRepository.GetByID(userID)
	if err != nil {
		return err
	}

	customer.WalletBalance += request.Amount
	if err := u.customerRepository.Update(customer); err != nil {
		return err
	}

	transaction := &models.Transaction{
		UserID:          userID,
		Amount:          request.Amount,
		PaymentMethodID: paymentMethod.ID,
		Type:            "topup",
	}

	if err := u.paymentRepository.CreateTransaction(transaction); err != nil {
		return err
	}

	return nil
}

func (u *PaymentUsecase) GetAllTransactions(userID uuid.UUID) ([]dto.TransactionDetail, error) {
	// _, err := u.customerRepository.GetByID(userID)
	// if err != nil {
	// 	return nil, err
	// }

	transactions, err := u.paymentRepository.GetAllTransactionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var transactionDetails []dto.TransactionDetail
	for _, tx := range transactions {
		paymentMethod, err := u.paymentRepository.GetPaymentMethodByID(tx.PaymentMethodID)
		if err != nil {
			return nil, err
		}
		transactionDetails = append(transactionDetails, dto.TransactionDetail{
			TransactionID: tx.ID,
			Amount:        tx.Amount,
			PaymentMethod: paymentMethod.Name,
			Type:          tx.Type,
			CreatedAt:     tx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return transactionDetails, nil
}

func (u *PaymentUsecase) PaidForFoodOrder(userID uuid.UUID, foodOrderID uuid.UUID) (*dto.PaymentSummary, error) {
	// ดึง order ที่ user คนนี้จ่าย
	order, err := u.paymentRepository.GetFoodOrderByID(foodOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get food order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("food order not found")
	}

	reservationID := order.ReservationID
	var summary *dto.PaymentSummary

	err = u.paymentRepository.RunInTransaction(func(tx *gorm.DB) error {
		fmt.Printf("[START] PaidForFoodOrder - user: %s, foodOrder: %s\n", userID, foodOrderID)

		member, err := u.paymentRepository.GetTableReservationMemberByCustomerID(reservationID, userID)
		if err != nil {
			return fmt.Errorf("failed to get reservation member: %w", err)
		}
		if member == nil {
			return fmt.Errorf("user is not a member of this reservation")
		}

		if member.Status == "completed" || member.Status == "paid" {
			return fmt.Errorf("user already confirmed payment")
		}

		if err := u.paymentRepository.UpdateTableReservationMemberStatus(member.ID, "paid_pending"); err != nil {
			return fmt.Errorf("failed to mark member paid_pending: %w", err)
		}
		fmt.Println("Member marked as paid_pending")

		members, err := u.paymentRepository.GetAllMembersByTableReservationID(reservationID)
		if err != nil {
			return fmt.Errorf("failed to get reservation members: %w", err)
		}

		totalMembers := len(members)
		pendingMembers := 0
		for _, m := range members {
			if m.Status == "paid" || m.Status == "completed" || m.Status == "paid_pending" {
				pendingMembers++
			}
		}
		fmt.Printf("Payment confirmations: %d/%d\n", pendingMembers, totalMembers)

		// ถ้ายังไม่ครบ ให้รอ
		if pendingMembers < totalMembers {
			fmt.Println("Waiting for other members to confirm payment...")
			summary = &dto.PaymentSummary{
				ReservationID: reservationID,
				FoodOrderID:   foodOrderID,
				TotalMembers:  totalMembers,
				PaidMembers:   pendingMembers,
			}
			return nil
		}

		// ✅ ครบทุกคนแล้ว หักเงินจริงทุกคน
		fmt.Println("All members confirmed. Proceeding with actual deduction...")

		// ดึง food orders ทั้งหมดใน reservation เดียวกัน
		orders, err := u.paymentRepository.GetAllFoodOrdersByReservationID(reservationID)
		if err != nil {
			return fmt.Errorf("failed to get all food orders: %w", err)
		}
		if len(orders) == 0 {
			return fmt.Errorf("no food orders found for this reservation")
		}

		// เตรียม payment method
		paymentMethods, err := u.paymentRepository.GetPaymentMethodsByType("paid")
		if err != nil || len(paymentMethods) == 0 {
			return fmt.Errorf("no valid payment method found for restaurant")
		}
		paymentMethod := paymentMethods[0]

		var totalPaid float64 = 0

		// หักเงินทุกคน
		for _, m := range members {
			var userTotal float64 = 0

			// หายอดรวมของ user แต่ละคนจากทุก order ใน reservation
			for _, ord := range orders {
				total, err := u.paymentRepository.GetTotalAmountForCustomerInOrder(ord.ID, m.CustomerID)
				if err != nil {
					return fmt.Errorf("failed to calculate total for member %v: %w", m.CustomerID, err)
				}
				userTotal += total
			}

			if userTotal <= 0 {
				fmt.Printf("Skip member %v: no orders found\n", m.CustomerID)
				continue
			}

			customer, err := u.customerRepository.GetByID(m.CustomerID)
			if err != nil {
				return fmt.Errorf("failed to get customer: %w", err)
			}

			if float64(customer.WalletBalance) < userTotal {
				return fmt.Errorf("insufficient balance for member %v: need %.2f, have %.2f",
					m.CustomerID, userTotal, customer.WalletBalance)
			}

			customer.WalletBalance -= float32(userTotal)
			if err := u.customerRepository.Update(customer); err != nil {
				return fmt.Errorf("failed to update wallet balance: %w", err)
			}

			tx := &models.Transaction{
				UserID:          m.CustomerID,
				PaymentMethodID: paymentMethod.ID,
				Amount:          float32(userTotal),
				Type:            "paid",
			}
			if err := u.paymentRepository.CreateTransaction(tx); err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			if err := u.paymentRepository.UpdateTableReservationMemberStatus(m.ID, "paid"); err != nil {
				return fmt.Errorf("failed to mark member paid: %w", err)
			}

			fmt.Printf("Deducted %.2f from member %v\n", userTotal, m.CustomerID)
			totalPaid += userTotal
		}

		// อัปเดตสถานะ reservation และ food orders
		if err := u.paymentRepository.UpdateTableReservationStatus(reservationID, "paid"); err != nil {
			return fmt.Errorf("failed to update reservation status: %w", err)
		}
		for _, ord := range orders {
			if err := u.paymentRepository.UpdateFoodOrderStatus(ord.ID, "paid"); err != nil {
				return fmt.Errorf("failed to update food order: %w", err)
			}
		}

		// รวมยอดทั้งหมด
		reservationTotal, err := u.paymentRepository.GetTotalAmountByReservationID(reservationID)
		if err != nil {
			return fmt.Errorf("failed to get total amount for reservation: %w", err)
		}

		restaurant, err := u.paymentRepository.GetRestaurantByFoodOrderID(orders[0].ID)
		if err != nil {
			return fmt.Errorf("failed to get restaurant: %w", err)
		}

		restaurant.WalletBalance += float32(reservationTotal)
		if err := u.restaurantRepository.Update(restaurant); err != nil {
			return fmt.Errorf("failed to update restaurant wallet: %w", err)
		}

		restTx := &models.Transaction{
			UserID:          restaurant.ID,
			PaymentMethodID: paymentMethod.ID,
			Amount:          float32(reservationTotal),
			Type:            "received",
		}
		if err := u.paymentRepository.CreateTransaction(restTx); err != nil {
			return fmt.Errorf("failed to create restaurant transaction: %w", err)
		}

		fmt.Printf("Restaurant wallet updated +%.2f (total: %.2f)\n", reservationTotal, restaurant.WalletBalance)

		summary = &dto.PaymentSummary{
			ReservationID: reservationID,
			FoodOrderID:   foodOrderID,
			TotalMembers:  totalMembers,
			PaidMembers:   totalMembers,
		}

		fmt.Println("[COMMIT] Transaction success")
		return nil
	})

	if err != nil {
		fmt.Printf("[ROLLBACK] Transaction failed: %v\n", err)
		return nil, err
	}

	fmt.Println("[END] PaidForFoodOrder completed successfully")
	return summary, nil
}

func (u *PaymentUsecase) GetWithdrawPaymentMethods(userID uuid.UUID) ([]dto.PaymentMethodDetail, error) {

	var paymentMethods []dto.PaymentMethodDetail
	methods, err := u.paymentRepository.GetPaymentMethodsByType("withdraw")
	if err != nil {
		return nil, err
	}

	for _, method := range methods {
		paymentMethods = append(paymentMethods, dto.PaymentMethodDetail{
			PaymentMethodID: method.ID,
			Name:            method.Name,
			// ImageURL:        method.ImageURL,
		})
	}

	return paymentMethods, nil
}

func (u *PaymentUsecase) WithdrawFromWallet(userID uuid.UUID, request *dto.WithdrawRequest) (*dto.WithdrawResponse, error) {
	restaurant, err := u.restaurantRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}

	withdrawMethods, err := u.paymentRepository.GetPaymentMethodsByType("withdraw")
	if err != nil || len(withdrawMethods) == 0 {
		return nil, fmt.Errorf("no valid payment method found for withdrawal")
	}

	var withdrawMethod *models.PaymentMethod
	for _, method := range withdrawMethods {
		if method.Name == request.BankName {
			withdrawMethod = &method
			break
		}
	}
	if withdrawMethod == nil {
		return nil, fmt.Errorf("invalid bank name")
	}

	var withdrawAmount float32 = request.WithdrawAmount
	if restaurant.WalletBalance < withdrawAmount {
		return nil, fmt.Errorf("insufficient wallet balance")
	}

	// Deduct the amount from the restaurant's wallet
	restaurant.WalletBalance -= withdrawAmount
	if err := u.restaurantRepository.Update(restaurant); err != nil {
		return nil, err
	}

	// Create a withdrawal transaction
	tx := &models.Transaction{
		UserID:          restaurant.ID,
		PaymentMethodID: withdrawMethod.ID,
		Amount:          -withdrawAmount,
		Type:            "withdraw",
	}
	if err := u.paymentRepository.CreateTransaction(tx); err != nil {
		return nil, err
	}

	return &dto.WithdrawResponse{
		RemainingBalance: restaurant.WalletBalance,
	}, nil
}
