// internal/order/usecase/get_detail_for_restaurant.go
package usecase

import (
	"backend/internal/order/dto"
	"context"
	"github.com/google/uuid"
)

func (u *orderUsecase) GetDetailForRestaurant(
	ctx context.Context,
	in GetDetailForRestaurantInput,
) (dto.OrderDetailForRestaurantResp, error) {

	raw, err := u.repo.GetOrderDetailForRestaurant(ctx, in.OrderID, in.RestaurantID)
	if err != nil {
		return dto.OrderDetailForRestaurantResp{}, err
	}

	// map options ตาม item
	itemOptsMap := map[uuid.UUID][]dto.OrderKitchenItemOption{}
	for _, op := range raw.Options {
		itemOptsMap[op.FoodOrderItemID] = append(itemOptsMap[op.FoodOrderItemID], dto.OrderKitchenItemOption{
			GroupName:  op.GroupName,
			OptionName: op.OptionName,
			Qty:        op.Qty,
		})
	}

	// map items
	items := make([]dto.OrderKitchenItem, 0, len(raw.Items))
	for _, it := range raw.Items {
		one := dto.OrderKitchenItem{
			OrderItemID:  it.ID,
			MenuItemID:   it.MenuItemID,
			MenuName:     it.MenuName,
			Quantity:     it.Quantity,
			UnitPrice:    it.UnitPrice,
			LineSubtotal: it.Subtotal,
			Note:         it.Note,
			Options:      itemOptsMap[it.ID],
		}
		items = append(items, one)
	}

	resp := dto.OrderDetailForRestaurantResp{
		OrderID:         raw.Order.ID,
		Status:          raw.Order.Status,
		OrderDate:       raw.Order.OrderDate,
		ExpectedReceive: raw.Order.ExpectedReceive,
		Note:            raw.Order.Note,
		Items:           items,
		TableLabel:      raw.TableLabel,
		TimeslotStart:   raw.TimeslotStart,
		TimeslotEnd:     raw.TimeslotEnd,
	}
	return resp, nil
}
