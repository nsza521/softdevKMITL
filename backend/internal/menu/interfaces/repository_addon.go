// repository_addon.go
package interfaces

import (
	models "backend/internal/db_model"
	"github.com/google/uuid"
)

type AddOnRepository interface {
	// Group
	CreateGroup(group *models.MenuAddOnGroup) error
	GetGroupsByRestaurant(restaurantID uuid.UUID) ([]models.MenuAddOnGroup, error)
	GetGroupByID(id uuid.UUID) (*models.MenuAddOnGroup, error)
	UpdateGroup(group *models.MenuAddOnGroup) error
	DeleteGroup(id uuid.UUID) error

	// Option
	CreateOption(opt *models.MenuAddOnOption) error
	GetOptionByID(id uuid.UUID) (*models.MenuAddOnOption, error)
	UpdateOption(opt *models.MenuAddOnOption) error
	DeleteOption(id uuid.UUID) error
}
