package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	// "backend/internal/db_model"
	"backend/internal/order/dto"
	"backend/internal/order/repository"
)

type OrderHistoryUsecase interface {
	GetServedHistoryForDay(ctx context.Context, restaurantID uuid.UUID, day time.Time) (*dto.OrderHistoryResponse, error)
}

type orderHistoryUsecase struct {
	repo repository.OrderHistoryRepository
}

func NewOrderHistoryUsecase(repo repository.OrderHistoryRepository) OrderHistoryUsecase {
	return &orderHistoryUsecase{repo: repo}
}

func (u *orderHistoryUsecase) GetServedHistoryForDay(
	ctx context.Context,
	restaurantID uuid.UUID,
	day time.Time,
) (*dto.OrderHistoryResponse, error) {

	orders, err := u.repo.ListServedOrdersByRestaurantAndDay(ctx, restaurantID, day)
	if err != nil {
		return nil, err
	}

	resp := dto.OrderHistoryResponse{
		Date:   day.Format(time.RFC3339), // หรือใช้ "2006-01-02"
		Orders: make([]dto.OrderHistoryOrder, 0, len(orders)),
	}

	for _, o := range orders {
		orderDTO := dto.OrderHistoryOrder{
			OrderID:     o.ID.String(),
			Status:      o.Status,
			Channel:     o.Channel,
			Note:        derefString(o.Note),
			TotalAmount: o.TotalAmount,
			OrderTime:   o.OrderDate.Format(time.RFC3339),
			Items:       []dto.OrderHistoryItem{},
		}

		for _, it := range o.Items {
			itemDTO := dto.OrderHistoryItem{
				MenuName:   it.MenuName,
				Quantity:   it.Quantity,
				UnitPrice:  it.UnitPrice,
				Subtotal:   it.Subtotal,
				Options:    []dto.OrderHistoryItemOption{},
			}

			for _, opt := range it.Options {
				itemDTO.Options = append(itemDTO.Options, dto.OrderHistoryItemOption{
					OptionName:  opt.OptionName,
					PriceDelta:  opt.PriceDelta,
				})
			}

			orderDTO.Items = append(orderDTO.Items, itemDTO)
		}

		resp.Orders = append(resp.Orders, orderDTO)
	}

	return &resp, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
