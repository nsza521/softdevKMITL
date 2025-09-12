package interfaces

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
}