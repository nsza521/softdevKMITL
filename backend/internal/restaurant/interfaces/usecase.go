package interfaces

import (
	"backend/internal/restaurant/dto"
	user "backend/internal/user/dto"
)

type RestaurantUsecase interface {
	Register(request *dto.RegisterRestaurantRequest) error
	Login(request *user.LoginRequest) (string, error)
}