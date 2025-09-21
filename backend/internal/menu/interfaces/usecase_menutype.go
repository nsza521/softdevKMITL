// internal/menu/interfaces/usecase_menutype.go
package interfaces

import (
	"context"

	models "backend/internal/db_model"
	menu "backend/internal/menu/dto"
	"github.com/google/uuid"
)

type MenuTypeUsecase interface {
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuType, error)
	Create(ctx context.Context, restaurantID uuid.UUID, in *menu.CreateMenuTypeRequest) (*models.MenuType, error)
	Update(ctx context.Context, typeID uuid.UUID, in *menu.UpdateMenuTypeRequest) (*models.MenuType, error)
	Delete(ctx context.Context, typeID uuid.UUID) error

	// เวอร์ชันตรวจ owner จาก JWT → typeID
	UpdateOwned(ctx context.Context, actorRestaurantID uuid.UUID, typeID uuid.UUID, in *menu.UpdateMenuTypeRequest) (*models.MenuType, error)
	DeleteOwned(ctx context.Context, actorRestaurantID uuid.UUID, typeID uuid.UUID) error
}
