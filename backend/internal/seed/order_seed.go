package seed

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/internal/db_model" // ปรับให้ตรง path จริงของ models package ของคุณ
)

// -----------------------------------------
// Minimal row structs for raw queries
// -----------------------------------------

type restaurantRow struct {
	ID       uuid.UUID
	Username string
}

type menuItemRow struct {
	ID           uuid.UUID
	Name         string
	Price        float64
	TimeTakenMin int `gorm:"column:time_taken"`
}

type addOnRow struct {
	OptID      uuid.UUID `gorm:"column:opt_id"`
	GroupID    uuid.UUID `gorm:"column:group_id"`
	GroupName  string    `gorm:"column:group_name"`
	OptionName string    `gorm:"column:option_name"`
	PriceDelta float64   `gorm:"column:price_delta"`
}

// -----------------------------------------
// Public entry
// -----------------------------------------

func RunSeedOrders(db *gorm.DB) error {
	const targetUsername = "restaurant_noodle"

	// 0) หา restaurant noodle
	rest, err := getNoodleRestaurant(db, targetUsername)
	if err != nil {
		return err
	}
	if rest == nil {
		return fmt.Errorf("cannot seed orders: restaurant %q not found", targetUsername)
	}

	// 1) preload candidate menu items (หลายจาน)
	menuItems, err := listMenuItemsForRestaurant(db, rest.ID, 5)
	if err != nil {
		return err
	}
	if len(menuItems) == 0 {
		return fmt.Errorf("cannot seed orders: noodle restaurant has no menu items")
	}

	// 2) preload add-ons (ทั้งร้าน)
	addOns, err := listAddOnsForRestaurant(db, rest.ID, 4)
	if err != nil {
		return err
	}
	// addOns ว่างได้

	// 3) เตรียมชุดออเดอร์ที่จะสร้าง
	//    ใช้ seedTag ต่างกันใน note เพื่อ idempotent แบบรายออเดอร์
	orderPlans := []orderPlan{
		{
			SeedTag:          "SEED_ORDER_SAMPLE_A",
			Channel:          "walk_in",
			Status:           "pending",
			ExpectedAfterMin: 15,
			ItemCount:        1,
			NoteText:         "SEED_ORDER_SAMPLE_A - โต๊ะ 3 ไม่ผัก",
		},
		{
			SeedTag:          "SEED_ORDER_SAMPLE_B",
			Channel:          "reservation",
			Status:           "preparing",
			ExpectedAfterMin: 25,
			ItemCount:        2,
			NoteText:         "SEED_ORDER_SAMPLE_B - ถึง 12:30",
		},
		{
			SeedTag:          "SEED_ORDER_SAMPLE_C",
			Channel:          "walk_in",
			Status:           "served",
			ExpectedAfterMin: 5,
			ItemCount:        1,
			NoteText:         "SEED_ORDER_SAMPLE_C - เอาเผ็ดน้อย",
		},
	}

	// 4) loop สร้างแต่ละออเดอร์แบบ transaction แยก
	for _, plan := range orderPlans {
		if err := seedOneOrderIfNotExists(db, rest.ID, menuItems, addOns, plan); err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------
// Plan/model helpers
// -----------------------------------------

// orderPlan อธิบายว่าจะสร้างออเดอร์หน้าตาแบบไหน
type orderPlan struct {
	SeedTag          string // ใช้ใน note เพื่อตรวจ idempotent
	Channel          string // walk_in | reservation | delivery
	Status           string // pending | preparing | served | paid | cancelled
	ExpectedAfterMin int    // เวลาเสิร์ฟคาดหวังกี่นาทีจาก now
	ItemCount        int    // อยากให้มีกี่จานในบิลนี้
	NoteText         string // note แสดงจริงใน DB
}

// snapshot menu item -> ใช้สร้าง FoodOrderItem
type menuItemSnapshot struct {
	MenuItemID   uuid.UUID
	MenuName     string
	UnitPrice    float64
	TimeTakenMin int
}

// snapshot addon -> ใช้สร้าง FoodOrderItemOption
type addOnSnapshot struct {
	AddOnOptionID uuid.UUID
	GroupID       uuid.UUID
	GroupName     string
	OptionName    string
	PriceDelta    float64
}

// -----------------------------------------
// Core seeding
// -----------------------------------------

func seedOneOrderIfNotExists(
	db *gorm.DB,
	restaurantID uuid.UUID,
	menuItems []menuItemRow,
	addOns []addOnRow,
	plan orderPlan,
) error {

	// ตรวจว่ามีอยู่แล้ว?
	var existing models.FoodOrder
	if err := db.
		// เราเช็คด้วย Note เพราะ FoodOrder ไม่มี RestaurantID ให้ filter.
		// Note เป็น pointer ในโมเดลจริงของคุณ, ใน DB เป็น column "note".
		// ของเราจะเซ็ต Note = plan.NoteText (มี seed tag ในตัว)
		Where("note = ?", plan.NoteText).
		First(&existing).Error; err == nil {
		return nil // มีอยู่แล้ว -> ข้าม
	}

	now := time.Now()
	expReceive := now.Add(time.Duration(plan.ExpectedAfterMin) * time.Minute)

	// เตรียม actor/dummy IDs ต่อ order
	customerID := uuid.New()
	// createdByUserID := uuid.New()
	var reservationID uuid.UUID
	if plan.Channel == "reservation" {
		reservationID = uuid.New() // จำลองว่ามีการจองโต๊ะ
	} else {
		reservationID = uuid.Nil // walk-in/delivery -> ไม่มี
	}

	// เตรียม FoodOrder struct
	orderID := uuid.New()
	orderNote := plan.NoteText
	order := models.FoodOrder{
		ID:              orderID,
		ReservationID:   reservationID,
		CustomerID:      customerID,
		// CreatedByUserID: createdByUserID,
		Status:          plan.Status,   // เช่น "pending" หรือ "preparing"
		Channel:         plan.Channel,  // "walk_in" / "reservation"
		OrderDate:       now,
		ExpectedReceive: &expReceive,
		TotalAmount:     0, // เราจะคำนวณจาก items แล้วอัปเดต ก่อน insert
		Note:            &orderNote,
		Items:           []models.FoodOrderItem{},
	}

	// เติม items ตาม ItemCount
	for i := 0; i < plan.ItemCount && i < len(menuItems); i++ {
		mirow := menuItems[i]

		itemID := uuid.New()
		qty := 1 // ถ้าอยาก random ก็ปรับตรงนี้ เช่น 1+i

		itemSubtotal := mirow.Price * float64(qty)

		// เลือก add-ons 0..len(addOns) (ไม่ต้องเยอะ เดี๋ยวมองไม่ออก)
		optionsForItem := []models.FoodOrderItemOption{}
		maxOpt := 0
		if len(addOns) > 0 {
			maxOpt = 1
			if len(addOns) > 1 && i%2 == 1 {
				maxOpt = 2 // ให้บางจานมี 2 ออปชัน
			}
		}
		for j := 0; j < maxOpt; j++ {
			opt := addOns[(i+j)%len(addOns)]
			itemSubtotal += opt.PriceDelta * float64(qty)

			optionsForItem = append(optionsForItem, models.FoodOrderItemOption{
				ID:              uuid.New(),
				FoodOrderItemID: itemID,
				AddOnOptionID:   opt.OptID,
				GroupID:         opt.GroupID,
				GroupName:       opt.GroupName,
				OptionName:      opt.OptionName,
				PriceDelta:      opt.PriceDelta,
				Qty:             1,
			})
		}

		itemNote := ""
		switch i {
		case 0:
			itemNote = "ไม่หวาน"
		case 1:
			itemNote = "เพิ่มเส้น"
		default:
			itemNote = "ใส่ทุกอย่าง"
		}

		order.Items = append(order.Items, models.FoodOrderItem{
			ID:           itemID,
			FoodOrderID:  orderID,
			CustomerID:   customerID,
			MenuItemID:   mirow.ID,
			MenuName:     mirow.Name,
			UnitPrice:    mirow.Price,
			TimeTakenMin: mirow.TimeTakenMin,
			Quantity:     qty,
			Subtotal:     itemSubtotal,
			Note:         strPtr(itemNote),
			Options:      optionsForItem,
		})

		order.TotalAmount += itemSubtotal
	}

	// insert
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("create seed order (%s): %w", plan.SeedTag, err)
		}
		return nil
	})
}

// -----------------------------------------
// Query helpers
// -----------------------------------------

func getNoodleRestaurant(db *gorm.DB, username string) (*restaurantRow, error) {
	var rest restaurantRow
	err := db.
		Table("restaurants").
		Select("id, username").
		Where("username = ?", username).
		Take(&rest).Error
	if err != nil {
		// ถ้าไม่เจอ -> nil, nil
		return nil, nil
	}
	return &rest, nil
}

// ดึงเมนูร้านนี้หลายรายการ (limit N) เพื่อนำไปสุ่มใส่ order
func listMenuItemsForRestaurant(db *gorm.DB, restaurantID uuid.UUID, limit int) ([]menuItemRow, error) {
	var rows []menuItemRow
	err := db.
		Table("menu_items").
		Select("id, name, price, time_taken").
		Where("restaurant_id = ?", restaurantID).
		Order("id").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list menu items for restaurant: %w", err)
	}
	return rows, nil
}

// ดึง add-on ของร้านนี้ limit N
func listAddOnsForRestaurant(db *gorm.DB, restaurantID uuid.UUID, limit int) ([]addOnRow, error) {
	var rows []addOnRow

	err := db.
		Table("menu_add_on_options AS o").
		Select(`
			o.id          AS opt_id,
			o.group_id    AS group_id,
			g.name        AS group_name,
			o.name        AS option_name,
			o.price_delta AS price_delta
		`).
		Joins("JOIN menu_add_on_groups AS g ON g.id = o.group_id").
		Where("g.restaurant_id = ?", restaurantID).
		Order("o.id").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		// ถ้า query fail (ยังไม่มี table แอดออน) -> return empty, nil
		return []addOnRow{}, nil
	}

	return rows, nil
}

// util เล็ก ๆ
func strPtr(s string) *string { return &s }
