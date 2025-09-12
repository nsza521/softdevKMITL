package http

import (
	// "github.com/gin-gonic/gin"

	"backend/internal/notifications/interfaces"
	// "backend/internal/notifications/dto"
)

type NotiHandler struct {
	notiUsecase interfaces.NotiUsecase
}

func NewNotiHandler(notiUsecase interfaces.NotiUsecase) interfaces.NotiHandler {
	return &NotiHandler{
		notiUsecase: notiUsecase,
	}
}