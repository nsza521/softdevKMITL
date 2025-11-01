// internal/menu/interfaces/repository_menutype.go
package interfaces

import (
	"context"

	"github.com/google/uuid"
	models "backend/internal/db_model"
)

type MenuTypeRepository interface {
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.MenuType, error)
	FindByName(ctx context.Context, restaurantID uuid.UUID, name string) (*models.MenuType, error)
	Create(ctx context.Context, mt *models.MenuType) error
	Update(ctx context.Context, mt *models.MenuType) error
	Delete(ctx context.Context, id uuid.UUID) error

	// เพิ่มสำหรับตรวจ owner
	GetMenuTypeRestaurantID(ctx context.Context, typeID uuid.UUID) (uuid.UUID, error)
}
