package seed

import (
	"fmt"
	"log"
	// "github.com/google/uuid"
	"gorm.io/gorm"
	"time"

	models "backend/internal/db_model"
)

func ptr[T any](v T) *T { return &v }

// SeedDev will insert sample Restaurant, Menu, MenuItem data
// Runs idempotently: if any menu exists, it won't double-insert.
func SeedDev(db *gorm.DB) error {
	// ถ้ามีเมนูอยู่แล้ว ไม่ seed ซ้ำ
	var count int64
	if err := db.Model(&models.Menu{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count menus: %w", err)
	}
	if count > 0 {
		return nil
	}

	// "15:04:05" is HH:mm:ss เป็น layout ของ Go
	// https://pkg.go.dev/time#Time.Format
	// https://pkg.go.dev/time#Parse
	open, _ := time.Parse("15:04:05", "08:00:00")
	close, _ := time.Parse("15:04:05", "20:00:00")

	rest := models.Restaurant{
		Username:      "kmitl_canteen",
		Password:      "hashedpassword123",
		Email:         "canteen@example.com",
		OpenTime:      models.OnlyTime{Time: open},
		CloseTime:     models.OnlyTime{Time: close},
		WalletBalance: 1000,
		ProfilePic:    ptr("restaurants/canteen.jpg"),
	}

	if err := db.Create(&rest).Error; err != nil {
		return fmt.Errorf("create restaurant: %w", err)
	}
	log.Printf("Created restaurant ID=%v", rest.ID)

	// เมนูหลัก + เครื่องดื่ม (อิง model เดิมของคุณ)
	mainMenu := models.Menu{
		Type:         "Main",
		RestaurantID: rest.ID,
	}
	drinkMenu := models.Menu{
		Type:         "Drink",
		RestaurantID: rest.ID,
	}
	if err := db.Create(&[]*models.Menu{&mainMenu, &drinkMenu}).Error; err != nil {
		return fmt.Errorf("create menus: %w", err)
	}
	log.Printf("mainMenu ID=%v, drinkMenu ID=%v", mainMenu.ID, drinkMenu.ID)


	items := []models.MenuItem{
		{
			Base:      models.Base{}, // ให้ GORM สร้าง ID เอง (ถ้า Base มี ID string และคุณอยากกำหนดเอง ก็ใส่ uuid.New())
			Name:      "Pad Thai",
			Price:     55,
			MenuID:    mainMenu.ID,
			MenuPic:   ptr("menus/pad-thai.jpg"), // เก็บเป็น key (จะอัปจริงใน MinIO ทีหลังได้)
			TimeTaken: 7,                         // นาที
		},
		{
			Base:      models.Base{},
			Name:      "Fried Rice",
			Price:     50,
			MenuID:    mainMenu.ID,
			MenuPic:   ptr("menus/fried-rice.jpg"),
			TimeTaken: 6,
		},
		{
			Base:      models.Base{},
			Name:      "Iced Latte",
			Price:     45,
			MenuID:    drinkMenu.ID,
			MenuPic:   ptr("menus/iced-latte.jpg"),
			TimeTaken: 3,
		},
		{
			Base:      models.Base{},
			Name:      "Thai Milk Tea",
			Price:     35,
			MenuID:    drinkMenu.ID,
			MenuPic:   ptr("menus/thai-milktea.jpg"),
			TimeTaken: 2,
		},
	}
	if err := db.Create(&items).Error; err != nil {
		return fmt.Errorf("create menu_items: %w", err)
	}

	return nil
}
