package interfaces

import (
	"time"
	"github.com/gin-gonic/gin"
	// "backend/internal/user/dto"
)

type UserHandler interface {
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
}

type UserRepository interface {

}

type UserUsecase interface {
	// Login(request *dto.LoginRequest) (string, error)
	Logout(token string, expiry time.Time) error
}