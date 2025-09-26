package repository

import (
	"backend/internal/db_model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FoodOrderRepository struct {
	db *gorm.DB
}

func NewFoodOrderRepository(db *gorm.DB) *FoodOrderRepository {
	return &FoodOrderRepository{db: db}
}

func (r *FoodOrderRepository) CreateOrderWithItems(ctx context.Context,
	order models.FoodOrder,
	items []models.FoodOrderItem,
	options []models.FoodOrderItemOption,
) (*models.FoodOrder, error) {
	var created models.FoodOrder

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].FoodOrderID = order.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// ต้อง map options ให้ตรงกับ item แต่ละตัว (ใน usecase ควรผูกมาให้แล้ว)
		for i := range options {
			if options[i].FoodOrderItemID == uuid.Nil {
				return gorm.ErrInvalidData
			}
		}
		if len(options) > 0 {
			if err := tx.Create(&options).Error; err != nil {
				return err
			}
		}

		created = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *FoodOrderRepository) GetOrderWithItems(ctx context.Context, orderID uuid.UUID) (*models.FoodOrder, error) {
	var order models.FoodOrder
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Options").
		First(&order, "id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
