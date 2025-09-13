package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	models "backend/internal/db_model"
	iface "backend/internal/menu/interfaces"
)

type menuUsecase struct{ repo iface.MenuRepository }

func NewMenuUsecase(r iface.MenuRepository) iface.MenuUsecase { return &menuUsecase{repo: r} }

func (u *menuUsecase) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]iface.MenuItemBrief, error) {
	items, err := u.repo.ListMenuByRestaurant(ctx, restaurantID)
	if err != nil { return nil, err }
	out := make([]iface.MenuItemBrief, 0, len(items))
	for _, m := range items {
		typeIDs := make([]uuid.UUID, 0, len(m.MenuTypes))
		types := make([]iface.MenuTypeBrief, 0, len(m.MenuTypes))
		for _, t := range m.MenuTypes {
			typeIDs = append(typeIDs, t.ID)
			types = append(types, iface.MenuTypeBrief{ID: t.ID, Type: t.Type})
		}
		out = append(out, iface.MenuItemBrief{
			ID: m.ID, Name: m.Name, Price: m.Price, MenuPic: m.MenuPic,
			TimeTaken: m.TimeTaken, Description: m.Description,
			MenuTypeIDs: typeIDs, Types: types,
		})
	}
	return out, nil
}

func (u *menuUsecase) CheckRestaurantExists(ctx context.Context, restaurantID uuid.UUID) error {
	return u.repo.RestaurantExists(ctx, restaurantID)
}

func (u *menuUsecase) CreateMenuItem(ctx context.Context, restaurantID uuid.UUID, in *iface.CreateMenuItemRequest) (*iface.MenuItemBrief, error) {
	if in == nil { return nil, errors.New("nil request") }
	if err := u.repo.VerifyMenuTypesBelongToRestaurant(ctx, restaurantID, in.MenuTypeIDs); err != nil {
		return nil, err
	}
	mi := &models.MenuItem{
		Name: in.Name, Price: in.Price, MenuPic: in.MenuPic,
		TimeTaken: in.TimeTaken, Description: in.Description,
	}
	if mi.TimeTaken == 0 { mi.TimeTaken = 1 }
	if err := u.repo.CreateMenuItem(ctx, mi); err != nil { return nil, err }
	if err := u.repo.AttachMenuTypes(ctx, mi.ID, in.MenuTypeIDs); err != nil { return nil, err }

	loaded, err := u.repo.LoadMenuItemWithTypes(ctx, mi.ID)
	if err != nil { return nil, err }
	resp := toBrief(loaded)
	return &resp, nil
}

func (u *menuUsecase) UpdateMenuItem(ctx context.Context, id uuid.UUID, in *iface.UpdateMenuItemRequest) (*iface.MenuItemBrief, error) {
	if in == nil {
		return nil, errors.New("nil request")
	}

	fields := map[string]any{}
	if in.Name != nil       { fields["name"]        = *in.Name }
	if in.Price != nil      { fields["price"]       = *in.Price }
	if in.TimeTaken != nil  { fields["time_taken"]  = *in.TimeTaken }
	if in.Description != nil{ fields["description"] = *in.Description }
	if in.MenuPic != nil    { fields["menu_pic"]    = *in.MenuPic }

	if len(fields) > 0 {
		if err := u.repo.UpdateMenuItem(ctx, id, fields); err != nil {
			return nil, err
		}
	}

	// ถ้ามีการส่ง menu_type_ids → แทนที่ใหม่ทั้งหมด
	if in.MenuTypeIDs != nil {
		if err := u.repo.ReplaceMenuTypes(ctx, id, *in.MenuTypeIDs); err != nil {
			return nil, err
		}
	}

	// โหลดข้อมูลล่าสุดกลับมาให้ response
	loaded, err := u.repo.LoadMenuItemWithTypes(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toBrief(loaded)
	return &resp, nil
}


func (u *menuUsecase) DeleteMenuItem(ctx context.Context, menuItemID uuid.UUID) error {
	return u.repo.DeleteMenuItem(ctx, menuItemID)
}

func toBrief(m *models.MenuItem) iface.MenuItemBrief {
	typeIDs := make([]uuid.UUID, 0, len(m.MenuTypes))
	types := make([]iface.MenuTypeBrief, 0, len(m.MenuTypes))
	for _, t := range m.MenuTypes {
		typeIDs = append(typeIDs, t.ID)
		types = append(types, iface.MenuTypeBrief{ID: t.ID, Type: t.Type})
	}
	return iface.MenuItemBrief{
		ID: m.ID, Name: m.Name, Price: m.Price, MenuPic: m.MenuPic,
		TimeTaken: m.TimeTaken, Description: m.Description,
		MenuTypeIDs: typeIDs, Types: types,
	}
}
