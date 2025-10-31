package usecase

import (
	"github.com/google/uuid"

	"backend/internal/payment/interfaces"
	"backend/internal/payment/dto"
)

type PaymentUsecase struct {
	paymentRepository interfaces.PaymentRepository
}

func NewPaymentUsecase(paymentRepository interfaces.PaymentRepository) interfaces.PaymentUsecase {
	return &PaymentUsecase{
		paymentRepository: paymentRepository,
	}
}

func (u *PaymentUsecase) GetTopupPaymentMethods(userID uuid.UUID) ([]dto.PaymentMethodDetail, error) {
	
	var paymentMethods []dto.PaymentMethodDetail

	methods, err := u.paymentRepository.GetTopupPaymentMethods()
	if err != nil {
		return nil, err
	}

	for _, method := range methods {
		paymentMethods = append(paymentMethods, dto.PaymentMethodDetail{
			PaymentMethodID: method.ID,
			Method:          method.Name,
			ImageURL:        method.ImageURL,
		})
	}

	return paymentMethods, nil
}
