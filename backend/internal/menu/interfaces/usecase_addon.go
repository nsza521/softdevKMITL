// usecase_addon.go
package interfaces

import (
	"backend/internal/menu/dto"
	"github.com/google/uuid"
)

type AddOnUsecase interface {
	// Group
	CreateGroup(restaurantID uuid.UUID, input menu.CreateAddOnGroupRequest) (menu.AddOnGroupResponse, error)
	ListGroups(restaurantID uuid.UUID) ([]menu.AddOnGroupResponse, error)
	UpdateGroup(id uuid.UUID, input menu.UpdateAddOnGroupRequest) error
	DeleteGroup(id uuid.UUID) error

	// Option
	GetOption(id uuid.UUID) (menu.AddOnOptionResponse, error)
	ListOptions(groupID uuid.UUID) ([]menu.AddOnOptionResponse, error)

	CreateOption(groupID uuid.UUID, input menu.CreateAddOnOptionRequest) (menu.AddOnOptionResponse, error)
	UpdateOption(id uuid.UUID, input menu.UpdateAddOnOptionRequest) error
	DeleteOption(id uuid.UUID) error
}
