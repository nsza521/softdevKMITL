package usecase

import (
	"fmt"
	"time"

	// "github.com/google/uuid"

	"backend/internal/user/dto"
	"backend/internal/user/interfaces"
	"backend/internal/utils"
	customerInterfaces "backend/internal/customer/interfaces"
	restaurantInterfaces "backend/internal/restaurant/interfaces"
)

type UserUsecase struct {
	userRepository   interfaces.UserRepository
	customerUsecase  customerInterfaces.CustomerUsecase
	restaurantUsecase restaurantInterfaces.RestaurantUsecase
}

func NewUserUsecase(userRepository interfaces.UserRepository, 
					customerUsecase customerInterfaces.CustomerUsecase, 
					restaurantUsecase restaurantInterfaces.RestaurantUsecase) interfaces.UserUsecase {
	return &UserUsecase{
		userRepository:  	userRepository,
		customerUsecase:    customerUsecase,
		restaurantUsecase:  restaurantUsecase,
	}
}

func (u *UserUsecase) Login(request *dto.LoginRequest) (string, error) {
	token, err := u.customerUsecase.Login(request)
	if err == nil {
		return token, nil
	}

	token, err = u.restaurantUsecase.Login(request)
	if err == nil {
		return token, nil
	}

	return "", fmt.Errorf("invalid username or password")

}


func (u *UserUsecase) Logout(token string, expiry time.Time) error {
	utils.BlacklistToken(token, expiry.Unix())
	return nil
}
