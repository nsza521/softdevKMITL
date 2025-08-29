package usecase

import (

	"backend/internal/customer/interfaces"
)

type CustomerUsecase struct {
	customerRepository interfaces.CustomerRepository
}

func NewCustomerUsecase(customerRepository interfaces.CustomerRepository) interfaces.CustomerUsecase {
	return &CustomerUsecase{
		customerRepository: customerRepository,
	}
}

func (u *CustomerUsecase) Register(username string) error {
	return nil
}

func (u *CustomerUsecase) Login(username string, password string) (string, error) {
	return "", nil
}
