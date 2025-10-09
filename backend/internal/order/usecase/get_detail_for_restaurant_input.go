package usecase

import "github.com/google/uuid"

type GetDetailForRestaurantInput struct {
	OrderID      uuid.UUID
	RestaurantID uuid.UUID
}
