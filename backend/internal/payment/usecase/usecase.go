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
			TransactionID:   tx.ID,
			Amount:          tx.Amount,
			PaymentMethod:   paymentMethod.Name,
			Type:            tx.Type,
			CreatedAt:       tx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return transactionDetails, nil
}

func (u *PaymentUsecase) PaidForFoodOrder(userID uuid.UUID, foodOrderID uuid.UUID) (*dto.PaymentSummary, error) {

    // 1. ดึง reservation ID จาก food order
    order, err := u.paymentRepository.GetFoodOrderByID(foodOrderID)
    if err != nil {
        return nil, fmt.Errorf("failed to get food order: %w", err)
    }

    reservationID := order.ReservationID

    var summary *dto.PaymentSummary

    err = u.paymentRepository.RunInTransaction(func(tx *gorm.DB) error {
        fmt.Printf("[START] PaidForFoodOrder - user: %s, foodOrder: %s\n", userID, foodOrderID)

        // Check if user is a member of the reservation
        member, err := u.paymentRepository.GetTableReservationMemberByCustomerID(reservationID, userID)
        if err != nil {
            return fmt.Errorf("failed to get reservation member: %w", err)
        }
        if member == nil {
            return fmt.Errorf("user is not a member of this reservation")
        }

        // Prevent duplicate payment confirmation
        if member.Status == "paid" || member.Status == "completed" {
            return fmt.Errorf("user already confirmed payment")
        }

        // Mark member as paid_pending
        if err := u.paymentRepository.UpdateTableReservationMemberStatus(member.ID, "paid_pending"); err != nil {
            return fmt.Errorf("failed to mark member paid_pending: %w", err)
        }
        fmt.Println("Member marked as paid_pending")

        // Check if all members have confirmed payment
        members, err := u.paymentRepository.GetAllMembersByTableReservationID(reservationID)
        if err != nil {
            return fmt.Errorf("failed to get reservation members: %w", err)
        }

        totalMembers := len(members)
        pendingMembers := 0
        for _, m := range members {
            if m.Status == "paid_pending" || m.Status == "paid" {
                pendingMembers++
            }
        }

        fmt.Printf("Payment confirmations: %d/%d\n", pendingMembers, totalMembers)

        // If not all members have confirmed, return pending status
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

        // 5. เมื่อทุกคนยืนยันครบ ให้เริ่มหักเงินจริงทั้งหมด
        fmt.Println("All members confirmed. Proceeding with actual deduction...")

        order, err := u.paymentRepository.GetFoodOrderByReservationID(reservationID)
        if err != nil {
            return fmt.Errorf("failed to get food order: %w", err)
        }
        if order == nil {
            return fmt.Errorf("no food order found for this reservation")
        }

        // 6. หักเงินจริงของสมาชิกแต่ละคนตามยอดรวม
        paymentMethods, err := u.paymentRepository.GetPaymentMethodsByType("paid")
        if err != nil || len(paymentMethods) == 0 {
            return fmt.Errorf("no valid payment method found to paid for restaurant")
        }

        paymentMethod := paymentMethods[0]

        var totalPaid float64 = 0
        for _, m := range members {
            userTotal, err := u.paymentRepository.GetTotalAmountForCustomerInOrder(order.ID, m.CustomerID)
            if err != nil {
                return fmt.Errorf("failed to calculate total for member %v: %w", m.CustomerID, err)
            }

            customer, err := u.customerRepository.GetByID(m.CustomerID)
            if err != nil {
                return fmt.Errorf("failed to get customer: %w", err)
            }

            if float64(customer.WalletBalance) < userTotal {
                return fmt.Errorf("insufficient balance for member %v: need %.2f, have %.2f",
                    m.CustomerID, userTotal, customer.WalletBalance)
            }

            // หักเงินจริง
            newBalance := customer.WalletBalance - float32(userTotal)
            customer.WalletBalance = newBalance
            if err := u.customerRepository.Update(customer); err != nil {
                return fmt.Errorf("failed to update wallet balance: %w", err)
            }

            // สร้าง transaction ของลูกค้า
            tx := &models.Transaction{
                UserID:          m.CustomerID,
                PaymentMethodID: paymentMethod.ID,
                Amount:          float32(userTotal),
                Type:            "paid",
            }
            if err := u.paymentRepository.CreateTransaction(tx); err != nil {
                return fmt.Errorf("failed to create transaction for member: %w", err)
            }

            // เปลี่ยนสถานะเป็น paid จริง
            if err := u.paymentRepository.UpdateTableReservationMemberStatus(m.ID, "paid"); err != nil {
                return fmt.Errorf("failed to mark member paid: %w", err)
            }

            fmt.Printf("Deducted %.2f from member %v\n", userTotal, m.CustomerID)
            totalPaid += userTotal
        }

        // 7. อัปเดตสถานะของ reservation และ food order
        if err := u.paymentRepository.UpdateTableReservationStatus(reservationID, "paid"); err != nil {
            return fmt.Errorf("failed to update reservation status: %w", err)
        }
        if err := u.paymentRepository.UpdateFoodOrderStatus(order.ID, "paid"); err != nil {
            return fmt.Errorf("failed to update food order: %w", err)
        }

        // 8. เพิ่มเงินให้ร้านค้า
        restaurant, err := u.paymentRepository.GetRestaurantByFoodOrderID(order.ID)
        if err != nil {
            return fmt.Errorf("failed to get restaurant: %w", err)
        }

        newRestBalance := restaurant.WalletBalance + float32(totalPaid)
        restaurant.WalletBalance = newRestBalance
        if err := u.restaurantRepository.Update(restaurant); err != nil {
            return fmt.Errorf("failed to update restaurant wallet: %w", err)
        }

        restTx := &models.Transaction{
            UserID:          restaurant.ID,
            PaymentMethodID: paymentMethod.ID,
            Amount:          float32(totalPaid),
            Type:            "received",
        }
        if err := u.paymentRepository.CreateTransaction(restTx); err != nil {
            return fmt.Errorf("failed to create restaurant transaction: %w", err)
        }

        fmt.Printf("Restaurant wallet updated +%.2f (total: %.2f)\n", totalPaid, newRestBalance)

        summary = &dto.PaymentSummary{
            ReservationID: reservationID,
            FoodOrderID:   order.ID,
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
