package http

import (
    "math"
    "net/http"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"

    "backend/internal/restaurant/interfaces"
)

type MenuDelivery struct {
    uc interfaces.MenuUsecase
}

func NewMenuDelivery(r *gin.Engine, uc interfaces.MenuUsecase) {
    h := &MenuDelivery{uc: uc}

    // เส้นทางหลัก: ดึงอาหารทั้งหมดของ “ร้าน”
    g := r.Group("/api/v1/restaurants/:restaurant_id")
    g.GET("/foods", h.GetFoods) // หรือใช้ชื่อ /menus/items ก็ได้ตามที่ทีมตกลง
}

func (d *MenuDelivery) GetFoods(c *gin.Context) {
    restaurantID := c.Param("restaurant_id")

    // รับ query params
    q := interfaces.MenuQuery{
        Q:        strings.TrimSpace(c.Query("q")),
        Type:     strings.TrimSpace(c.Query("type")), // filter ที่ menus.type
        Sort:     c.DefaultQuery("sort", "created_at"),
        Order:    c.DefaultQuery("order", "desc"),
        Page:     atoiDefault(c.Query("page"), 1),
        PageSize: atoiDefault(c.Query("page_size"), 20),
    }
    if v := c.Query("min_price"); v != "" {
        if f, err := strconv.ParseFloat(v, 64); err == nil {
            q.MinPrice = &f
        }
    }
    if v := c.Query("max_price"); v != "" {
        if f, err := strconv.ParseFloat(v, 64); err == nil {
            q.MaxPrice = &f
        }
    }

    items, total, err := d.uc.GetRestaurantFoods(c, restaurantID, q)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    totalPages := int(math.Ceil(float64(total) / float64(q.PageSize)))
    c.JSON(http.StatusOK, gin.H{
        "items":       items,
        "page":        q.Page,
        "page_size":   q.PageSize,
        "total":       total,
        "total_pages": totalPages,
    })
}

func atoiDefault(s string, def int) int {
    if s == "" {
        return def
    }
    n, err := strconv.Atoi(s)
    if err != nil || n <= 0 {
        return def
    }
    return n
}
