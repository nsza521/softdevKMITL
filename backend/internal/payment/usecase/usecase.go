package usecase

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/google/uuid"

	"backend/internal/payment/dto"
	"backend/internal/payment/interfaces"
	"backend/internal/db_model"
	customerInterfaces "backend/internal/customer/interfaces"
	restaurantInterfaces "backend/internal/restaurant/interfaces"
)

type PaymentUsecase struct {
	paymentRepository interfaces.PaymentRepository
	customerRepository customerInterfaces.CustomerRepository
	restaurantRepository restaurantInterfaces.RestaurantRepository
}

func NewPaymentUsecase(paymentRepository interfaces.PaymentRepository, 
	customerRepository customerInterfaces.CustomerRepository, 
	restaurantRepository restaurantInterfaces.RestaurantRepository,
	) interfaces.PaymentUsecase {
	return &PaymentUsecase{
		paymentRepository: paymentRepository,
		customerRepository: customerRepository,
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
	if paymentMethod.Type != "topup" && paymentMethod.Type != "all" {
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
	_, err := u.customerRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}

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
			TransactionID:   tx.ID,
			Amount:          tx.Amount,
			PaymentMethod:   paymentMethod.Name,
			Type:            tx.Type,
			CreatedAt:       tx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return transactionDetails, nil
}

func (u *PaymentUsecase) PaidForFoodOrder(userID uuid.UUID, reservationID uuid.UUID) (*dto.PaymentSummary, error) {
    var summary *dto.PaymentSummary

    err := u.paymentRepository.RunInTransaction(func(tx *gorm.DB) error {
        fmt.Printf("💳 [START] PaidForFoodOrder - user: %s, reservation: %s\n", userID, reservationID)

        // 1️⃣ ตรวจสอบว่า user เป็นสมาชิกของ reservation นี้ไหม
        member, err := u.paymentRepository.GetTableReservationMemberByCustomerID(reservationID, userID)
        if err != nil {
            fmt.Printf("❌ Failed to get reservation member: %v\n", err)
            return err
        }
        if member == nil {
            return fmt.Errorf("user is not a member of this reservation")
        }
        if member.Status == "paid" {
            return fmt.Errorf("user has already paid")
        }

        // 2️⃣ โหลด FoodOrder ของโต๊ะนี้
        order, err := u.paymentRepository.GetFoodOrderByReservationID(reservationID)
        if err != nil {
            fmt.Printf("❌ Failed to get food order: %v\n", err)
            return err
        }
        if order == nil {
            return fmt.Errorf("no food order found for this reservation")
        }

        // 3️⃣ รวมยอดอาหารของลูกค้าคนนี้
        userTotal, err := u.paymentRepository.GetTotalAmountForCustomerInOrder(order.ID, userID)
        if err != nil {
            fmt.Printf("❌ Failed to calculate user's total: %v\n", err)
            return err
        }
        if userTotal <= 0 {
            return fmt.Errorf("no food items found for this user in the order")
        }
        fmt.Printf("🧾 User total: %.2f\n", userTotal)

        // 4️⃣ ตรวจสอบยอดใน Wallet
        customer, err := u.customerRepository.GetByID(userID)
        if err != nil {
            fmt.Printf("❌ Failed to get customer: %v\n", err)
            return err
        }
        fmt.Printf("👛 Current wallet: %.2f\n", customer.WalletBalance)
        if float64(customer.WalletBalance) < userTotal {
            return fmt.Errorf("insufficient balance: need %.2f, have %.2f", userTotal, customer.WalletBalance)
        }

        // 5️⃣ หักเงินออกจาก Wallet
        newBalance := customer.WalletBalance - float32(userTotal)
		customer.WalletBalance = newBalance
        if err := u.customerRepository.Update(customer); err != nil {
            fmt.Printf("❌ Failed to update wallet balance: %v\n", err)
            return err
        }
        fmt.Printf("✅ Wallet updated: %.2f → %.2f\n", customer.WalletBalance, newBalance)

        // 6️⃣ สร้าง Transaction
        transaction := &models.Transaction{
            UserID:          userID,
            PaymentMethodID: uuid.Nil,
            Amount:          float32(userTotal),
            Type:            "paid",
        }
        if err := u.paymentRepository.CreateTransaction(transaction); err != nil {
            fmt.Printf("❌ Failed to create transaction: %v\n", err)
            return err
        }
        fmt.Println("🧾 Transaction created successfully")

        // 7️⃣ อัปเดต member เป็น “paid”
        if err := u.paymentRepository.UpdateTableReservationMemberStatus(member.ID, "paid"); err != nil {
            fmt.Printf("❌ Failed to update reservation member: %v\n", err)
            return err
        }
        fmt.Println("✅ Member marked as paid")

        // 8️⃣ ตรวจสอบว่าทุกคนจ่ายครบหรือยัง
        members, err := u.paymentRepository.GetAllMembersByTableReservationID(reservationID)
        if err != nil {
            fmt.Printf("❌ Failed to get reservation members: %v\n", err)
            return err
        }

        totalMembers := len(members)
        paidMembers := 0
        for _, m := range members {
            if m.CustomerID == userID {
                m.Status = "paid"
            }
            if m.Status == "paid" {
                paidMembers++
            }
        }
        fmt.Printf("👥 Paid members: %d/%d\n", paidMembers, totalMembers)

        // 9️⃣ ถ้าทุกคนจ่ายครบ
        if paidMembers == totalMembers {
            if err := u.paymentRepository.UpdateTableReservationStatus(reservationID, "paid"); err != nil {
                fmt.Printf("❌ Failed to update reservation: %v\n", err)
                return err
            }
            if err := u.paymentRepository.UpdateFoodOrderStatus(order.ID, "paid"); err != nil {
                fmt.Printf("❌ Failed to update food order: %v\n", err)
                return err
            }
            fmt.Println("🎉 All members paid! Reservation and order marked as 'paid'")
        }

        summary = &dto.PaymentSummary{
            ReservationID: reservationID,
            FoodOrderID:   order.ID,
            TotalMembers:  totalMembers,
            PaidMembers:   paidMembers,
        }

        fmt.Println("✅ [COMMIT] Transaction success")
        return nil
    })

    if err != nil {
        fmt.Printf("🚨 [ROLLBACK] Transaction failed: %v\n", err)
        return nil, err
    }

    fmt.Println("💰 [END] PaidForFoodOrder completed successfully")
    return summary, nil
}
