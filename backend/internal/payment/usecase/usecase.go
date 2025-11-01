package usecase

import (
	"fmt"
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