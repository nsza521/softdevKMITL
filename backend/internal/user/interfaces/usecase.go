package interfaces

import (
	"time"
	// "github.com/google/uuid"

	"backend/internal/user/dto"
)

type CustomerLogin interface {
    Login(*dto.LoginRequest) (string, error)
}

type RestaurantLogin interface {
    Login(*dto.LoginRequest) (string, error)
}

type UserUsecase interface {
	Login(request *dto.LoginRequest) (string, error)
	Logout(token string, expiry time.Time) error
}