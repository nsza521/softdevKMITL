// internal/menu/usecase/usecase.go
package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	interfaces "backend/internal/menu/interfaces"
)

type menuUsecase struct{ 
	repo interfaces.MenuRepository 
	minioClient *minio.Client
}

func NewMenuUsecase(r interfaces.MenuRepository, minioClient *minio.Client) interfaces.MenuUsecase { 
	return &menuUsecase{
		repo: r, 
		minioClient: minioClient,
	} 
}

func (u *menuUsecase) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]interfaces.MenuItemBrief, error) {

	// check ว่า restaurantID มีอยู่จริงไหม

	items, err := u.repo.ListMenuByRestaurant(ctx, restaurantID)
	if err != nil { return nil, err }

	out := make([]interfaces.MenuItemBrief, 0, len(items))
	for _, m := range items {
		typeIDs := make([]uuid.UUID, 0, len(m.MenuTypes))
		types   := make([]interfaces.MenuTypeBrief, 0, len(m.MenuTypes))
		for _, t := range m.MenuTypes {
			typeIDs = append(typeIDs, t.ID)
			types   = append(types, interfaces.MenuTypeBrief{ID: t.ID, Type: t.Type})
		}
		out = append(out, interfaces.MenuItemBrief{
			ID: m.ID, Name: m.Name, Price: m.Price, MenuPic: m.MenuPic,
			TimeTaken: m.TimeTaken, Description: m.Description,
			MenuTypeIDs: typeIDs,
			Types:       types, // 👈 ติด tag รายละเอียดมาด้วย
		})
	}
	return out, nil
}

func (u *menuUsecase) CheckRestaurantExists(ctx context.Context, restaurantID uuid.UUID) error {
	return u.repo.RestaurantExists(ctx, restaurantID)
}
