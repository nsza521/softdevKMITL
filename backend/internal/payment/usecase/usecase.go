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

	methods, err := u.paymentRepository.GetPaymentMethods("topup")
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

func (u *PaymentUsecase) TopupToWallet(userID uuid.UUID, amount float32, paymentMethodID uuid.UUID) error {
	paymentMethod, err := u.paymentRepository.GetPaymentMethodByID(paymentMethodID)
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

	customer.WalletBalance += amount
	if err := u.customerRepository.Update(customer); err != nil {
		return err
	}

	transaction := &models.Transaction{
		UserID:          userID,
		Amount:          amount,
		PaymentMethodID: paymentMethodID,
		Type:            "topup",
	}

	if err := u.paymentRepository.CreateTransaction(transaction); err != nil {
		return err
	}

	return nil
}

func (u *PaymentUsecase) CreateTopupTransaction(userID uuid.UUID, amount float32, paymentMethodID uuid.UUID) error {
	// Implementation for creating a top-up transaction
	return nil
}