package repository

import (
	"gorm.io/gorm"
)

type FoodOrderRepository struct {
	db *gorm.DB
}

func NewFoodOrderRepository(db *gorm.DB) *FoodOrderRepository {
	return &FoodOrderRepository{
		db: db,
	}
}