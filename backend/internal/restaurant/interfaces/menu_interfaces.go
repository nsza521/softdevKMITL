package interfaces

import (
    "context"
    "backend/internal/db_model"
)

type MenuQuery struct {
    Q        string  // keyword ที่จะค้นหาในชื่ออาหาร
    Type     string  // filter ที่ Menu.Type (เช่น "Drink", "Main")
    MinPrice *float64
    MaxPrice *float64
    Page     int
    PageSize int
    Sort     string // name|price|created_at
    Order    string // asc|desc
}

type MenuRepository interface {
    // ดึงอาหารทั้งหมดของร้านผ่านการ join: menus -> menu_items
    FindItemsByRestaurant(ctx context.Context, restaurantID string, q MenuQuery) (items []models.MenuItem, total int64, err error)
}

type MenuUsecase interface {
    GetRestaurantFoods(ctx context.Context, restaurantID string, q MenuQuery) (items []models.MenuItem, total int64, err error)
}
