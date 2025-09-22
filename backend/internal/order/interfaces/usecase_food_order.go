package interfaces

import (
	"backend/internal/order/dto"
	"github.com/google/uuid"
	"context"
)

type FoodOrderUsecase interface {
	Create(ctx context.Context, req dto.CreateFoodOrderReq, currentUser uuid.UUID) (dto.CreateFoodOrderResp, error)
	GetDetail(ctx context.Context, orderID uuid.UUID) (dto.FoodOrderResponse, error)

	// ใหม่:
	AppendItems(ctx context.Context, orderID uuid.UUID, req dto.AppendItemsReq, currentUser uuid.UUID) (dto.AppendItemsResp, error)
	RemoveItem(ctx context.Context, orderID uuid.UUID, itemID uuid.UUID, currentUser uuid.UUID) (dto.RemoveItemResp, error)
	AttachCustomer(ctx context.Context, orderID uuid.UUID, req dto.AttachCustomerReq, currentUser uuid.UUID) (dto.ListCustomersResp, error)
	ListCustomers(ctx context.Context, orderID uuid.UUID, currentUser uuid.UUID) (dto.ListCustomersResp, error)
}
