package usecase

import (

	"backend/internal/payment/interfaces"
)

type PaymentUsecase struct {
	paymentRepository interfaces.PaymentRepository
}

func NewPaymentUsecase(paymentRepository interfaces.PaymentRepository) interfaces.PaymentUsecase {
	return &PaymentUsecase{
		paymentRepository: paymentRepository,
	}
}

