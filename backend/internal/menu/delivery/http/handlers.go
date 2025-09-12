package http

import (
	// "github.com/gin-gonic/gin"

	"backend/internal/menu/interfaces"
	// "backend/internal/menu/dto"
)

type MenuHandler struct {
	menuUsecase interfaces.MenuUsecase
}

func NewMenuHandler(menuUsecase interfaces.MenuUsecase) interfaces.MenuHandler {
	return &MenuHandler{
		menuUsecase: menuUsecase,
	}
}