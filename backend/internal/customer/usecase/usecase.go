package usecase

import (
	"fmt"

	"backend/internal/customer/dto"
	"backend/internal/customer/interfaces"
	"backend/internal/utils"
	"backend/internal/db_model"
	user "backend/internal/user/dto"
)

type CustomerUsecase struct {
	customerRepository interfaces.CustomerRepository
}

func NewCustomerUsecase(customerRepository interfaces.CustomerRepository) interfaces.CustomerUsecase {
	return &CustomerUsecase{
		customerRepository: customerRepository,
	}
}

func (u *CustomerUsecase) Register(request *dto.RegisterCustomerRequest) error {
	
	// Check if customer exists
	exists, err := u.customerRepository.IsCustomerExists(request.Username, request.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("customer already exists")
	}

	// Validate email format
	if !utils.IsValidEmail(request.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Check password strength
	if !utils.IsStrongPassword(request.Password) {
		return fmt.Errorf("password is not strong enough")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	// Create new customer
	customer := models.Customer{
		Username:     request.Username,
		Email:        request.Email,
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Password:     hashedPassword,
	}

	return  u.customerRepository.Create(&customer)
}

func (u *CustomerUsecase) Login(request *user.LoginRequest) (string, error) {
	customer, err := u.customerRepository.GetByUsername(request.Username)
	if err != nil {
		return "", err
	}

	// hashedPassword, err := utils.HashPassword(request.Password)
	// if err != nil {
	// 	return "", err
	// }

	// Check password
	err = utils.VerifyPassword(request.Password, customer.Password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(customer.ID, customer.Username, "customer")
	if err != nil {
		return "", err
	}
	return token, nil

}

func (u *CustomerUsecase) Logout() error {
	// Implement logout logic if needed (e.g., token invalidation)
	return nil
}
