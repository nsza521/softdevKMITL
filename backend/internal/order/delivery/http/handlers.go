package http

import (
	// "github.com/gin-gonic/gin"

	"backend/internal/order/interfaces"
	// "backend/internal/customer/dto"
)

type FoodOrderHandler struct {
	foodOrderUsecase interfaces.FoodOrderUsecase
}

func NewFoodOrderHandler(foodOrderUsecase interfaces.FoodOrderUsecase) interfaces.FoodOrderHandler {
	return &FoodOrderHandler{
		foodOrderUsecase: foodOrderUsecase,
	}
}
