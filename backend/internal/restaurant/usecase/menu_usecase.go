package usecase

import (
    "context"
    "backend/internal/restaurant/interfaces"
)

type menuUsecase struct {
    repo interfaces.MenuRepository
}

func NewMenuUsecase(r interfaces.MenuRepository) *menuUsecase {
    return &menuUsecase{repo: r}
}

func (u *menuUsecase) GetRestaurantFoods(ctx context.Context, restaurantID string, q interfaces.MenuQuery) ([]interface{}, int64, error) {
    // สามารถเติม validation เพิ่มได้ เช่น ตรวจรูปแบบ UUID, limit, ฯลฯ
    items, total, err := u.repo.FindItemsByRestaurant(ctx, restaurantID, q)
    if err != nil {
        return nil, 0, err
    }
    // ถ้าต้อง map ออกเป็น DTO ก็ทำจุดนี้
    out := make([]interface{}, 0, len(items))
    for _, it := range items {
        out = append(out, it) // ส่ง model ตรง ๆ ไปก่อนตามที่ขอ "ไม่แก้ model"
    }
    return out, total, nil
}
