package usecase

import (
	"context"
	"fmt"

	"backend/internal/db_model"
	"backend/internal/order/dto"
	"backend/internal/order/repository"
)

type QueueUsecase interface {
	GetQueue(ctx context.Context, actorUserID string, actorRole string, restaurantID string) (dto.QueueResponse, error)
}

type queueUsecase struct {
	repo repository.QueueRepository
}

func NewQueueUsecase(repo repository.QueueRepository) QueueUsecase {
	return &queueUsecase{repo: repo}
}

func (u *queueUsecase) GetQueue(ctx context.Context, actorUserID string, actorRole string, restaurantID string) (dto.QueueResponse, error) {
	
	orders, err := u.repo.ListPendingOrdersByRestaurant(ctx, restaurantID)
	if err != nil {
		return dto.QueueResponse{}, err
	}

	resp := dto.QueueResponse{Orders: make([]dto.QueueOrderDTO, 0)}

	for _, fo := range orders {
		filtered := make([]models.FoodOrderItem, 0)
		for _, it := range fo.Items {
			ok, err := u.repo.CountMenuItemBelongsToRestaurant(it.MenuItemID.String(), restaurantID)
			if err != nil {
				return dto.QueueResponse{}, err
			}
			if ok {
				filtered = append(filtered, it)
			}
		}

		if len(filtered) == 0 {
			continue
		}
		fo.Items = filtered
		resp.Orders = append(resp.Orders, mapFoodOrderToDTO(fo))
	}

	return resp, nil
}

// map model → DTO
func mapFoodOrderToDTO(fo models.FoodOrder) dto.QueueOrderDTO {
	out := dto.QueueOrderDTO{
		ID:              fo.ID.String(),
		Status:          fo.Status,
		Channel:         fo.Channel,
		OrderDate:       fo.OrderDate,
		ExpectedReceive: fo.ExpectedReceive,
		TotalAmount:     fo.TotalAmount,
		Note:            fo.Note,
		Items:           make([]dto.QueueItemDTO, 0, len(fo.Items)),
	}
	fmt.Printf("Mapping FoodOrder ID %s with %d items\n", fo.ID.String(), len(fo.Items))

	for _, it := range fo.Items {
		fmt.Printf("  Item ID %s: MenuName=%s, Quantity=%d, UnitPrice=%.2f, MenuPic=%s\n", it.ID.String(), it.MenuName, it.Quantity, it.UnitPrice, it.MenuPic)
	}

	for _, it := range fo.Items {
		itemDTO := dto.QueueItemDTO{
			ID:           it.ID.String(),
			MenuName:     it.MenuName,
			Quantity:     it.Quantity,
			UnitPrice:    it.UnitPrice,
			Subtotal:     it.Subtotal,
			TimeTakenMin: it.TimeTakenMin,
			Note:         it.Note,
			Options:      make([]dto.QueueItemOptionDTO, 0, len(it.Options)),
			MenuPic:      it.MenuPic,
		}

		for _, opt := range it.Options {
			itemDTO.Options = append(itemDTO.Options, dto.QueueItemOptionDTO{
				OptionName: opt.OptionName,
				GroupName:  opt.GroupName,
				PriceDelta: opt.PriceDelta,
				Qty:        opt.Qty,
			})
		}

		out.Items = append(out.Items, itemDTO)
	}

	return out
}
