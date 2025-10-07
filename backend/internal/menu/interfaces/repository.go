package interfaces

import (
	"context"

	"github.com/google/uuid"
	models "backend/internal/db_model"
)

type MenuRepository interface {
	ListMenuByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuItem, error)
	RestaurantExists(ctx context.Context, restaurantID uuid.UUID) error

	// GetItemWithTypesAndAddOns retrieves a MenuItem by its ID, preloading its associated MenuTypes and AddOnGroups with their Options.
	GetItemWithTypesAndAddOns(itemID uuid.UUID) (*models.MenuItem, error)

	CreateMenuItem(ctx context.Context, mi *models.MenuItem) error
	UpdateMenuItem(ctx context.Context, id uuid.UUID, fields map[string]any) error
	DeleteMenuItem(ctx context.Context, id uuid.UUID) error

	AttachMenuTypes(ctx context.Context, itemID uuid.UUID, typeIDs []uuid.UUID) error
	ReplaceMenuTypes(ctx context.Context, itemID uuid.UUID, typeIDs []uuid.UUID) error

	VerifyMenuTypesBelongToRestaurant(ctx context.Context, restaurantID uuid.UUID, typeIDs []uuid.UUID) error
	LoadMenuItemWithTypes(ctx context.Context, id uuid.UUID) (*models.MenuItem, error)

	GetMenuItemByID(ctx context.Context, id uuid.UUID) (*models.MenuItem, error)
	
}
