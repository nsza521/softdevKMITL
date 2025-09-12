package interfaces

import (
	// "backend/internal/db_model"
	"backend/internal/customer/dto"
	user "backend/internal/user/dto"
)

type CustomerUsecase interface {
	Register(request *dto.RegisterCustomerRequest) error
	Login(request *user.LoginRequest) (string, error)
	Logout() error
}