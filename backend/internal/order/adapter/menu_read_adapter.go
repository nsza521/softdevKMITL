// internal/order/adapter/menu_read_adapter.go
package adapter

import (
	"context"

	"github.com/google/uuid"

	// === เปลี่ยน path ให้ตรงโปรเจกต์คุณ ===
	menudto "backend/internal/menu/dto"

	orderuc "backend/internal/order/usecase"
)

// ใช้ interface บางเบา เพื่อหลีกเลี่ยงการดึง usecase ทั้งแพ็กเกจมาเป็น dependency
// คุณสามารถเปลี่ยนให้เป็น usecase จริงของโมดูล menu ได้ ถ้ามีเมธอดชื่อ GetDetail ตรงกัน
type MenuQuery interface {
	GetDetail(itemID uuid.UUID) (menudto.MenuItemDetailResponse, error)
}

// MenuReadAdapter ทำหน้าที่แปลงผลจาก GetDetail → orderuc.MenuDetail
type MenuReadAdapter struct {
	Menu MenuQuery
}

func NewMenuReadAdapter(menu MenuQuery) *MenuReadAdapter {
	return &MenuReadAdapter{Menu: menu}
}

// Order usecase เรียกผ่านเมธอดนี้
func (a *MenuReadAdapter) GetMenuDetail(ctx context.Context, id uuid.UUID) (*orderuc.MenuDetail, error) {
	// NOTE: เมธอดจริงของคุณไม่มี context → ไม่ใช้ ctx ที่นี่
	res, err := a.Menu.GetDetail(id)
	if err != nil {
		return nil, err
	}

	// --- map -> orderuc.MenuDetail ---
	md := &orderuc.MenuDetail{
		ID:           res.ID,
		RestaurantID: res.RestaurantID,
		Name:         res.Name,
		Price:        res.Price,
		MenuPic:      res.MenuPic,
		TimeTaken:    res.TimeTaken,
		Description:  &res.Description,
	}

	// types
	for _, t := range res.Types {
		md.Types = append(md.Types, struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		}{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	// addons
	for _, g := range res.AddOns {
		gg := orderuc.AddOnGroup{
			ID:        g.ID,
			Name:      g.Name,
			Required:  g.Required,
			MinSelect: func() int {
				if g.MinSelect != nil {
					return *g.MinSelect
				}
				return 0
			}(),
			MaxSelect: func() int {
				if g.MaxSelect != nil {
					return *g.MaxSelect
				}
				return 0
			}(),
			AllowQty:  g.AllowQty,
			From:      g.From,
		}
		for _, op := range g.Options {
			gg.Options = append(gg.Options, orderuc.AddOption{
				ID:         op.ID,
				Name:       op.Name,
				PriceDelta: op.PriceDelta,
				IsDefault:  op.IsDefault,
				MaxQty:     op.MaxQty,
			})
		}
		md.AddOns = append(md.AddOns, gg)
	}

	return md, nil
}
