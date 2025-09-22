package usecase

import (
	"backend/internal/order/dto"
	orderif "backend/internal/order/interfaces"
	dbm "backend/internal/db_model"
	"context"

	"github.com/google/uuid"
)

type foodOrderUsecase struct {
	repo orderif.FoodOrderRepository
}

func NewFoodOrderUsecase(repo orderif.FoodOrderRepository) orderif.FoodOrderUsecase {
	return &foodOrderUsecase{repo: repo}
}

func (u *foodOrderUsecase) Create(ctx context.Context, req dto.CreateFoodOrderReq, currentUser uuid.UUID) (dto.CreateFoodOrderResp, error) {
	var items   []dbm.FoodOrderItem
	var options []dbm.FoodOrderItemOption
	var total   float64

	for _, it := range req.Items {
		// TODO: โหลดราคาเมนู/ออปชันจริงจาก repo เมนู
		unit := 100.0
		var addonSubtotal float64

		item := dbm.FoodOrderItem{
			MenuItemID:          it.MenuItemID,
			Quantity:            it.Quantity,
			UnitPrice:           unit,
			CreatedByCustomerID: &currentUser,
		}

		for _, o := range it.Options {
			priceDelta := 10.0 // TODO: ดึงจริงจาก MenuAddOnOption
			addonSubtotal += priceDelta * float64(o.Qty)

			options = append(options, dbm.FoodOrderItemOption{
				// NOTE: FoodOrderItemID จะถูกเซ็ตใน repo หลังสร้าง item
				AddOnOptionID: o.OptionID,
				Qty:           o.Qty,
				PriceDelta:    priceDelta,
			})
		}

		item.AddOnSubtotal = addonSubtotal
		item.LineTotal     = (unit + addonSubtotal) * float64(it.Quantity)
		total += item.LineTotal

		items = append(items, item)
	}

	order := dbm.FoodOrder{
		RestaurantID:  req.RestaurantID,
		ReservationID: req.ReservationID,
		Channel:       req.Channel,
		ExpectedTime:  req.ExpectedTime,
		Status:        "pending",
		Notes:         req.Notes,
		TotalAmount:   total,
	}

	created, err := u.repo.CreateOrderWithItems(ctx, order, items, options)
	if err != nil {
		return dto.CreateFoodOrderResp{}, err
	}

	return dto.CreateFoodOrderResp{
		OrderID:      created.ID,
		Status:       created.Status,
		TotalAmount:  created.TotalAmount,
		ExpectedTime: created.ExpectedTime,
	}, nil
}

func (u *foodOrderUsecase) GetDetail(ctx context.Context, orderID uuid.UUID) (dto.FoodOrderResponse, error) {
	order, err := u.repo.GetOrderWithItems(ctx, orderID)
	if err != nil {
		return dto.FoodOrderResponse{}, err
	}

	items := make([]dto.FoodOrderItemDTO, 0, len(order.Items))
	for _, it := range order.Items {
		opts := make([]dto.FoodOrderItemOptionDTO, 0, len(it.Options))
		for _, o := range it.Options {
			opts = append(opts, dto.FoodOrderItemOptionDTO{
				ID:         o.ID,
				OptionID:   o.AddOnOptionID,
				Qty:        o.Qty,
				PriceDelta: o.PriceDelta,
			})
		}
		items = append(items, dto.FoodOrderItemDTO{
			ID:                it.ID,
			MenuItemID:        it.MenuItemID,
			Quantity:          it.Quantity,
			UnitPrice:         it.UnitPrice,
			AddOnSubtotal:     it.AddOnSubtotal,
			LineTotal:         it.LineTotal,
			CreatedByCustomer: it.CreatedByCustomerID,
			Options:           opts,
		})
	}

	return dto.FoodOrderResponse{
		ID:            order.ID,
		RestaurantID:  order.RestaurantID,
		ReservationID: order.ReservationID,
		ExpectedTime:  order.ExpectedTime,
		Status:        order.Status,
		Notes:         order.Notes,
		TotalAmount:   order.TotalAmount,
		Channel:       order.Channel,
		Items:         items,
	}, nil
}

func (u *foodOrderUsecase) AppendItems(ctx context.Context, orderID uuid.UUID, req dto.AppendItemsReq, currentUser uuid.UUID) (dto.AppendItemsResp, error) {
	// TODO: ลอจิกจริง = ตรวจสิทธิ์ reservation, โหลดราคาเมนู/option, คำนวณยอด, อัปเดตใน TX
	// ตอนนี้ return แบบ stub
	order, err := u.repo.GetOrderWithItems(ctx, orderID)
	if err != nil { return dto.AppendItemsResp{}, err }
	return dto.AppendItemsResp{
		OrderID:     order.ID,
		AddedCount:  len(req.Items),
		TotalAmount: order.TotalAmount, // (ของจริงต้องคำนวณแล้วอัปเดต)
	}, nil
}

func (u *foodOrderUsecase) RemoveItem(ctx context.Context, orderID uuid.UUID, itemID uuid.UUID, currentUser uuid.UUID) (dto.RemoveItemResp, error) {
	// TODO: ลบใน TX + คำนวณยอดใหม่
	// ตอนนี้ stub
	_, err := u.repo.GetOrderWithItems(ctx, orderID)
	if err != nil { return dto.RemoveItemResp{}, err }
	return dto.RemoveItemResp{
		OrderID:     orderID,
		RemovedItem: itemID,
		TotalAmount: 0, // ของจริงต้องคำนวณใหม่
	}, nil
}

func (u *foodOrderUsecase) AttachCustomer(ctx context.Context, orderID uuid.UUID, req dto.AttachCustomerReq, currentUser uuid.UUID) (dto.ListCustomersResp, error) {
	// TODO: ตรวจว่า customer อยู่ใน reservation เดียวกัน, upsert ลง pivot, แล้วรีเทิร์นรายการลูกค้าทั้งหมด
	// ตอนนี้ stub
	order, err := u.repo.GetOrderWithItems(ctx, orderID)
	if err != nil { return dto.ListCustomersResp{}, err }
	return dto.ListCustomersResp{
		OrderID:   order.ID,
		Customers: []dto.OrderCustomerDTO{{CustomerID: req.CustomerID, Role: "contributor"}},
	}, nil
}

func (u *foodOrderUsecase) ListCustomers(ctx context.Context, orderID uuid.UUID, currentUser uuid.UUID) (dto.ListCustomersResp, error) {
	// TODO: ดึงจาก pivot table จริง (FoodOrderCustomers)
	// ตอนนี้ stub
	_, err := u.repo.GetOrderWithItems(ctx, orderID)
	if err != nil { return dto.ListCustomersResp{}, err }
	return dto.ListCustomersResp{
		OrderID:   orderID,
		Customers: []dto.OrderCustomerDTO{}, // ยังว่าง
	}, nil
}


