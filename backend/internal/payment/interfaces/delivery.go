package interfaces

import (
	"github.com/gin-gonic/gin"
)

type CustomerHandler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
}