package http

import (
	// "github.com/gin-gonic/gin"

	"backend/internal/table/interfaces"
	// "backend/internal/customer/dto"
)

type TableHandler struct {
	tableUsecase interfaces.TableUsecase
}

func NewTableHandler(tableUsecase interfaces.TableUsecase) interfaces.TableHandler {
	return &TableHandler{
		tableUsecase: tableUsecase,
	}
}