package interfaces

import (
	"context"

	"github.com/google/uuid"
	models "backend/internal/db_model"
)

type MenuRepository interface {
	ListMenuByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuItem, error)
	RestaurantExists(ctx context.Context, restaurantID uuid.UUID) error
}
