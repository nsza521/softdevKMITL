package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	 models "backend/internal/db_model" // <— ให้ alias เป็น models ให้ตรงกับการใช้งานด้านล่าง
	"backend/internal/order/dto"
	"backend/internal/order/repository"
)

/***************  Menu Read-Model (match JSON ของคุณ)  ***************/
type MenuDetail struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	MenuPic      *string   `json:"menu_pic"`
	TimeTaken    int       `json:"time_taken"`
	Description  *string   `json:"description"`
	Types        []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"types"`
	AddOns []AddOnGroup `json:"addons"`
}

type AddOnGroup struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Required  bool        `json:"required"`
	MinSelect int         `json:"min_select"`
	MaxSelect int         `json:"max_select"`
	AllowQty  bool        `json:"allow_qty"`
	From      string      `json:"from"` // "type" | "item"
	Options   []AddOption `json:"options"`
}
type AddOption struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceDelta float64   `json:"price_delta"`
	IsDefault  bool      `json:"is_default"`
	MaxQty     *int      `json:"max_qty"`
}

/***************  Port to read menu detail  ***************/
type MenuReadService interface {
	GetMenuDetail(ctx context.Context, menuItemID uuid.UUID) (*MenuDetail, error)
}

/***************  Usecase  ***************/
type OrderUsecase interface {
	// เวอร์ชันใหม่: reservation_id อยู่ใน body (optional)
	Create(ctx context.Context, req dto.CreateFoodOrderReq, currentCustomer uuid.UUID) (dto.CreateFoodOrderResp, error)
	GetDetailForRestaurant(ctx context.Context, input GetDetailForRestaurantInput) (dto.OrderDetailForRestaurantResp, error)
}

type orderUsecase struct {
	repo  repository.OrderRepository
	menu  MenuReadService
	nowFn func() time.Time
}

func NewOrderUsecase(repo repository.OrderRepository, menu MenuReadService) OrderUsecase {
	return &orderUsecase{repo: repo, menu: menu, nowFn: time.Now}
}

func (u *orderUsecase) Create(ctx context.Context, req dto.CreateFoodOrderReq, currentCustomer uuid.UUID) (dto.CreateFoodOrderResp, error) {
	if len(req.Items) == 0 {
		return dto.CreateFoodOrderResp{}, errors.New("no items")
	}

	fmt.Printf("CreateFoodOrderReq: %+v\n", req)
	// 1) โหลด reservation เฉพาะกรณีมี reservation_id ใน body
	var rsv *repository.Reservation
	if req.ReservationID != nil {
		rr, err := u.repo.LoadReservationForCustomer(ctx, *req.ReservationID, currentCustomer)
		if err != nil {
			return dto.CreateFoodOrderResp{}, err
		}
		rsv = rr
		fmt.Printf("Loaded reservation22: %+v\n", rsv)
	}else{
		fmt.Printf("No reservation ID provided, proceeding without reservation.\n")
	}

	

	order := &models.FoodOrder{
		ID: uuid.New(),
		// ReservationID: จะถูกเซ็ตใน repo ถ้ามี rsv (หรือจะเซ็ตเองที่นี่ก็ได้)
		ReservationID: rsv.ReservationID,
		// CustomerID: currentCustomer,
		Status:     "pending",
		OrderDate:  u.nowFn(),
		Note:       req.Note,
	}

	var orderItems []models.FoodOrderItem
	var orderOpts []models.FoodOrderItemOption
	var total float64

	// ป้องกัน cross-tenant: รวม restaurant ของทุกรายการ
	var restID *uuid.UUID

	for _, it := range req.Items {
		if it.Quantity <= 0 {
			return dto.CreateFoodOrderResp{}, errors.New("quantity must be >= 1")
		}
		// 2) โหลด menu detail (already merged type∪item)
		detail, err := u.menu.GetMenuDetail(ctx, it.MenuItemID)
		if err != nil {
			return dto.CreateFoodOrderResp{}, err
		}

		// (Guard) ทุก item ต้องเป็นร้านเดียวกัน
		if restID == nil {
			rid := detail.RestaurantID
			restID = &rid
		} else if *restID != detail.RestaurantID {
			return dto.CreateFoodOrderResp{}, errors.New("items belong to multiple restaurants")
		}

		// 3) ทำ map group/option
		type groupCtx struct {
			group AddOnGroup
			opts  map[uuid.UUID]AddOption
			picks int // จำนวน selections ใน group (นับเป็นรายการ ไม่ใช่ qty)
		}
		groups := map[uuid.UUID]*groupCtx{}
		for _, g := range detail.AddOns {
			m := map[uuid.UUID]AddOption{}
			for _, op := range g.Options {
				m[op.ID] = op
			}
			cp := g
			groups[g.ID] = &groupCtx{group: cp, opts: m}
		}

		// 4) validate selections + คำนวณ addon subtotal
		addonSubtotal := 0.0
		var itemOpts []models.FoodOrderItemOption

		for _, sel := range it.Selections {
			gc, ok := groups[sel.GroupID]
			if !ok {
				return dto.CreateFoodOrderResp{}, fmt.Errorf("invalid selection: unknown group (%s)", sel.GroupID)
			}
			op, ok := gc.opts[sel.OptionID]
			if !ok {
				return dto.CreateFoodOrderResp{}, fmt.Errorf("invalid selection: option not in group (%s)", sel.OptionID)
			}
			qty := sel.Qty
			if !gc.group.AllowQty {
				qty = 1
			} else {
				if qty <= 0 {
					return dto.CreateFoodOrderResp{}, errors.New("qty must be > 0 for allow_qty group")
				}
				if op.MaxQty != nil && qty > *op.MaxQty {
					return dto.CreateFoodOrderResp{}, errors.New("qty exceeds max_qty")
				}
			}
			gc.picks++

			addonSubtotal += float64(qty) * op.PriceDelta
			itemOpts = append(itemOpts, models.FoodOrderItemOption{
				ID:            uuid.New(),
				AddOnOptionID: op.ID,
				GroupID:       gc.group.ID,
				GroupName:     gc.group.Name,
				OptionName:    op.Name,
				PriceDelta:    op.PriceDelta,
				Qty:           qty,
			})
		}

		// 5) ตรวจ min/max per group
		for _, gc := range groups {
			if gc.group.Required && gc.picks == 0 {
				return dto.CreateFoodOrderResp{}, fmt.Errorf("required group not selected: %s", gc.group.Name)
			}
			if gc.picks > 0 {
				if gc.picks < gc.group.MinSelect {
					return dto.CreateFoodOrderResp{}, fmt.Errorf("min_select not met for group: %s", gc.group.Name)
				}
				if gc.group.MaxSelect > 0 && gc.picks > gc.group.MaxSelect {
					return dto.CreateFoodOrderResp{}, fmt.Errorf("max_select exceeded for group: %s", gc.group.Name)
				}
			}
		}

		lineBase := detail.Price
		lineSubtotal := (lineBase + addonSubtotal) * float64(it.Quantity)

		item := models.FoodOrderItem{
			ID:           uuid.New(),
			MenuItemID:   detail.ID,
			CustomerID:   currentCustomer,
			MenuName:     detail.Name,
			UnitPrice:    detail.Price,
			TimeTakenMin: detail.TimeTaken,
			Quantity:     it.Quantity,
			Subtotal:     lineSubtotal,
			Note:         it.Note,
		}
		for i := range itemOpts {
			itemOpts[i].FoodOrderItemID = item.ID
		}
		orderItems = append(orderItems, item)
		orderOpts = append(orderOpts, itemOpts...)
		total += lineSubtotal
	}

	order.TotalAmount = total

	// 6) persist (transaction) — repo จะ set ReservationID ให้ถ้ามี rsv
	if err := u.repo.CreateOrderTx(ctx, order, orderItems, orderOpts); err != nil {
		return dto.CreateFoodOrderResp{}, err
	}

	return dto.CreateFoodOrderResp{
		OrderID:     order.ID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}, nil
}
