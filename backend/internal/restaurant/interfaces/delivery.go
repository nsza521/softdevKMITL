package interfaces

import (
	"github.com/gin-gonic/gin"
)

type RestaurantHandler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	GetAll() gin.HandlerFunc
}