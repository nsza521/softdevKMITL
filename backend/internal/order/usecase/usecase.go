package usecase

import (

	"backend/internal/order/interfaces"
)

type FoodOrderUsecase struct {
	foodOrderRepository interfaces.FoodOrderRepository
}

func NewFoodOrderUsecase(foodOrderRepository interfaces.FoodOrderRepository) interfaces.FoodOrderUsecase {
	return &FoodOrderUsecase{
		foodOrderRepository: foodOrderRepository,
	}
}

