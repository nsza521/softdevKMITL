package repository

import (
    "context"
    "math"
    "strings"

    "gorm.io/gorm"

    "backend/internal/db_model"
    "backend/internal/restaurant/interfaces"
)

type menuRepository struct {
    db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *menuRepository {
    return &menuRepository{db}
}

func (r *menuRepository) FindItemsByRestaurant(ctx context.Context, restaurantID string, q interfaces.MenuQuery) ([]models.MenuItem, int64, error) {
    // ตารางตาม default naming ของ GORM: menus, menu_items
    base := r.db.WithContext(ctx).
        Model(&models.MenuItem{}).
        Joins("JOIN menus ON menus.id = menu_items.menu_id").
        Where("menus.restaurant_id = ?", restaurantID)

    if s := strings.TrimSpace(q.Q); s != "" {
        like := "%" + s + "%"
        base = base.Where("menu_items.name LIKE ?", like)
    }
    if t := strings.TrimSpace(q.Type); t != "" {
        base = base.Where("menus.type = ?", t) // filter ที่ field Type ของ Menu
    }
    if q.MinPrice != nil {
        base = base.Where("menu_items.price >= ?", *q.MinPrice)
    }
    if q.MaxPrice != nil {
        base = base.Where("menu_items.price <= ?", *q.MaxPrice)
    }

    // นับจำนวนทั้งหมดก่อน
    var total int64
    if err := base.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // จัดหน้า + sort
    page := q.Page
    if page <= 0 {
        page = 1
    }
    size := q.PageSize
    if size <= 0 || size > 100 {
        size = 20
    }

    sort := strings.ToLower(q.Sort)
    switch sort {
    case "name", "price", "created_at":
    default:
        sort = "created_at"
    }
    order := strings.ToUpper(q.Order)
    if order != "ASC" {
        order = "DESC"
    }

    var items []models.MenuItem
    err := base.
        Order("menu_items." + sort + " " + order).
        Offset((page-1)*size).
        Limit(size).
        Find(&items).Error
    if err != nil {
        return nil, 0, err
    }
    _ = math.Ceil // ป้องกัน go vet เตือนถ้าไม่ได้ใช้; เผื่อคุณอยากคำนวณ total_pages ฝั่ง handler

    return items, total, nil
}
