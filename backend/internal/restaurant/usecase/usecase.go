package usecase

import (
	"fmt"

	"backend/internal/restaurant/interfaces"
	"backend/internal/restaurant/dto"
	"backend/internal/utils"
	"backend/internal/db_model"
	user "backend/internal/user/dto"
)

type RestaurantUsecase struct {
	restaurantRepository interfaces.RestaurantRepository
}

func NewRestaurantUsecase(restaurantRepository interfaces.RestaurantRepository) interfaces.RestaurantUsecase {
	return &RestaurantUsecase{
		restaurantRepository: restaurantRepository,
	}
}

func (u *RestaurantUsecase) Register(request *dto.RegisterRestaurantRequest) error {

	// Check if restaurant exists
	exists, err := u.restaurantRepository.IsRestaurantExists(request.Username, request.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("restaurant already exists")
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

	// Create new restaurant
	restaurant := models.Restaurant{
		Username:     request.Username,
		Email:        request.Email,
		Password:     hashedPassword,
	}
	createdRestaurant, err := u.restaurantRepository.Create(&restaurant)
	if err != nil {
		return err
	}

	// Create bank account
	bankAccount := models.BankAccount{
		UserID:        createdRestaurant.ID,
		BankName:      request.BankName,
		AccountNumber: request.AccountNumber,
		AccountName:   request.AccountName,
	}
	err = u.restaurantRepository.CreateBankAccount(&bankAccount)
	if err != nil {
		return err
	}

	return nil
}

func (u *RestaurantUsecase) Login(request *user.LoginRequest) (string, error) {
	restaurant, err := u.restaurantRepository.GetByUsername(request.Username)
	if err != nil {
		return "", err
	}

	// Check password
	err = utils.VerifyPassword(request.Password, restaurant.Password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(restaurant.ID, restaurant.Username, "restaurant")
	if err != nil {
		return "", err
	}
	return token, nil

}

