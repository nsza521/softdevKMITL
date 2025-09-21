// internal/menu/usecase/menutype_uc.go
package usecase

import (
	"context"
	"errors"
	"strings"

	models "backend/internal/db_model"
	menu "backend/internal/menu/dto"
	"github.com/google/uuid"

	iface "backend/internal/menu/interfaces"
)

type menuTypeUC struct {
	repo iface.MenuTypeRepository
}

func NewMenuTypeUsecase(r iface.MenuTypeRepository) iface.MenuTypeUsecase {
	return &menuTypeUC{repo: r}
}

func (u *menuTypeUC) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuType, error) {
	return u.repo.ListByRestaurant(ctx, restaurantID)
}

func (u *menuTypeUC) Create(ctx context.Context, restaurantID uuid.UUID, in *menu.CreateMenuTypeRequest) (*models.MenuType, error) {
	name := strings.TrimSpace(in.Type)
	if name == "" {
		return nil, errors.New("type required")
	}

	existingTypes, err := u.repo.FindByName(ctx, restaurantID, name)
	if err != nil {
        return nil, err
    }
    if existingTypes != nil {
        return nil, errors.New("menu type already exists in this restaurant")
    }
	mt := &models.MenuType{
		Base:         models.Base{ID: uuid.New()},
		Type:         name,
		RestaurantID: restaurantID,
	}
	if err := u.repo.Create(ctx, mt); err != nil {
		return nil, err
	}
	return mt, nil
}

func (u *menuTypeUC) Update(ctx context.Context, typeID uuid.UUID, in *menu.UpdateMenuTypeRequest) (*models.MenuType, error) {
	mt, err := u.repo.GetByID(ctx, typeID)
	if err != nil {
		return nil, err
	}
	if mt == nil {
		return nil, errors.New("menu type not found")
	}

	if in.Type != nil {
		name := strings.TrimSpace(*in.Type)
		if name == "" {
			return nil, errors.New("type cannot be empty")
		}
		mt.Type = name
	}

	if err := u.repo.Update(ctx, mt); err != nil {
		return nil, err
	}
	return mt, nil
}

func (u *menuTypeUC) Delete(ctx context.Context, typeID uuid.UUID) error {
	return u.repo.Delete(ctx, typeID)
}


// ---- Owned versions ----
func (u *menuTypeUC) UpdateOwned(ctx context.Context, actorRestaurantID uuid.UUID, typeID uuid.UUID, in *menu.UpdateMenuTypeRequest) (*models.MenuType, error) {
	ownerID, err := u.repo.GetMenuTypeRestaurantID(ctx, typeID)
	if err != nil {
		return nil, err
	}
	if ownerID == uuid.Nil {
		return nil, errors.New("menu type not found")
	}
	if actorRestaurantID == uuid.Nil || actorRestaurantID != ownerID {
		return nil, errors.New("forbidden: not restaurant owner")
	}
	return u.Update(ctx, typeID, in)
}

func (u *menuTypeUC) DeleteOwned(ctx context.Context, actorRestaurantID uuid.UUID, typeID uuid.UUID) error {
	ownerID, err := u.repo.GetMenuTypeRestaurantID(ctx, typeID)
	if err != nil {
		return err
	}
	if ownerID == uuid.Nil {
		return errors.New("menu type not found")
	}
	if actorRestaurantID == uuid.Nil || actorRestaurantID != ownerID {
		return errors.New("forbidden: not restaurant owner")
	}
	return u.Delete(ctx, typeID)
}
