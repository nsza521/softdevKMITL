package main

import (
	"fmt"
	"time"

	// "log"
	// "os"
	models "backend/internal/db_model"
	"backend/internal/utils"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

func InitAllSeedData(db *gorm.DB)  error {

	err := seedCustomers(db)
	if err != nil {
		return fmt.Errorf("error seeding customers: %v", err)
	}

	err = seedRestaurants(db)
	if err != nil {
		return fmt.Errorf("error seeding restaurants: %v", err)
	}

	if err := seedMenuTypesAndItems(db); err != nil {
		return fmt.Errorf("error seeding menu types & items: %v", err)
	}

	err = seedTableTimeslots(db)
	if err != nil {
		return fmt.Errorf("error seeding tables and time slots: %v", err)
	}

	return nil
}

func seedCustomers(db *gorm.DB) error {

	for i := 1; i <= 10; i++ {
		username := fmt.Sprintf("customer%02d", i)
		email := fmt.Sprintf("customer%02d@example.com", i)

		// Check if customer already exists
		var count int64
		db.Model(&models.Customer{}).Where("username = ? OR email = ?", username, email).Count(&count)
		if count > 0 {
			continue
		}
		
		// Hash password
		hashedPassword, err := utils.HashPassword("SecureP@ssw0rd")
		if err != nil {
			return err
		}

		// Create new customer
		customer := models.Customer{
			Username: username,
			Email:    email,
			Password: hashedPassword,
			FirstName: fmt.Sprintf("FirstName%02d", i),
			LastName:  fmt.Sprintf("LastName%02d", i),
		}
		if err := db.Create(&customer).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedRestaurants(db *gorm.DB) error {
	for i := 1; i <= 10; i++ {
		username := fmt.Sprintf("restaurant%02d", i)
		email := fmt.Sprintf("restaurant%02d@example.com", i)

		// Check if restaurant already exists
		var count int64
		db.Model(&models.Restaurant{}).Where("username = ? OR email = ?", username, email).Count(&count)
		if count > 0 {
			continue
		}
		
		// Hash password
		hashedPassword, err := utils.HashPassword("SecureP@ssw0rd")
		if err != nil {
			return err
		}

		// Create new restaurant
		restaurant := models.Restaurant{
			Username: username,
			Email:    email,
			Password: hashedPassword,
		}
		if err := db.Create(&restaurant).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedMenuTypesAndItems(db *gorm.DB) error {
	var restaurants []models.Restaurant
	if err := db.Find(&restaurants).Error; err != nil {
		return fmt.Errorf("list restaurants: %w", err)
	}
	if len(restaurants) == 0 {
		return fmt.Errorf("no restaurants to seed menu for")
	}

	for _, rest := range restaurants {
		if err := seedMenuForRestaurant(db, rest.ID); err != nil {
			return err
		}
	}
	return nil
}

func seedMenuForRestaurant(db *gorm.DB, restaurantID uuid.UUID) error {
	// 1) เตรียม/หา MenuType ของร้าน ถ้ามีอยู่แล้วจะไม่สร้างซ้ำ
	categoryNames := []string{
		"อาหารจานเดียว",
		"ของหวาน",
		"เครื่องดื่ม",
		"ท็อปปิ้ง",
	}
	typeIDMap := make(map[string]uuid.UUID)

	for _, name := range categoryNames {
		mt, err := getOrCreateMenuType(db, restaurantID, name)
		if err != nil {
			return err
		}
		typeIDMap[name] = mt.ID
	}

	// 2) สร้าง/หา MenuItem
	type itemSpec struct {
		Name        string
		Price       float64
		TimeTaken   int
		Description string
		Types       []string // จะผูกกับชื่อ MenuType ด้านบน
	}

	items := []itemSpec{
		{"ข้าวกะเพรา", 65, 5, "เผ็ดกลาง", []string{"อาหารจานเดียว", "ท็อปปิ้ง"}},
		{"ผัดไทยกุ้งสด", 85, 7, "เส้นเหนียวนุ่ม", []string{"อาหารจานเดียว"}},
		{"ชาเย็น", 40, 2, "หวานหอม", []string{"เครื่องดื่ม"}},
		{"ไอศกรีมกะทิ", 45, 2, "หอมมัน", []string{"ของหวาน"}},
		{"ชาดำเย็น", 35, 2, "หวานกลาง", []string{"เครื่องดื่ม"}},
	}

	for _, it := range items {
		menuItem, created, err := getOrCreateMenuItem(db, it.Name, it.Price, it.TimeTaken, it.Description)
		if err != nil {
			return err
		}

		// 3) ผูกความสัมพันธ์ผ่าน MenuTag
		//    - ถ้าพึ่งสร้างใหม่ → ใส่ tag ตาม spec
		//    - ถ้าเคยมีอยู่แล้ว → ตรวจว่ามี tag ครบหรือยัง ขาดค่อยเติม
		if created {
			if err := attachTypes(db, menuItem.ID, it.Types, typeIDMap); err != nil {
				return err
			}
		} else {
			if err := ensureTypes(db, menuItem.ID, it.Types, typeIDMap); err != nil {
				return err
			}
		}
	}

	return nil
}

func getOrCreateMenuType(db *gorm.DB, restaurantID uuid.UUID, typeName string) (*models.MenuType, error) {
	var mt models.MenuType
	if err := db.
		Where("`type` = ? AND restaurant_id = ?", typeName, restaurantID).
		First(&mt).Error; err == nil {
		return &mt, nil
	}

	mt = models.MenuType{
		Type:         typeName,
		RestaurantID: restaurantID,
	}
	if err := db.Create(&mt).Error; err != nil {
		return nil, fmt.Errorf("create menutype %q: %w", typeName, err)
	}
	return &mt, nil
}

func getOrCreateMenuItem(db *gorm.DB, name string, price float64, timeTaken int, desc string) (*models.MenuItem, bool, error) {
	var mi models.MenuItem
	if err := db.Where("name = ?", name).First(&mi).Error; err == nil {
		// อัปเดตราคา/desc เบา ๆ ให้ทันสมัย (ถ้าต้องการ)
		mi.Price = price
		mi.Description = desc
		if timeTaken > 0 {
			mi.TimeTaken = timeTaken
		}
		if err := db.Model(&models.MenuItem{Base: mi.Base}).Updates(map[string]any{
			"price":       mi.Price,
			"time_taken":  mi.TimeTaken,
			"description": mi.Description,
		}).Error; err != nil {
			return nil, false, fmt.Errorf("update menuitem %q: %w", name, err)
		}
		return &mi, false, nil
	}

	mi = models.MenuItem{
		Name:        name,
		Price:       price,
		TimeTaken:   timeTaken,
		Description: desc,
	}
	if mi.TimeTaken == 0 {
		mi.TimeTaken = 1
	}
	if err := db.Create(&mi).Error; err != nil {
		return nil, false, fmt.Errorf("create menuitem %q: %w", name, err)
	}
	return &mi, true, nil
}

func attachTypes(db *gorm.DB, menuItemID uuid.UUID, typeNames []string, typeIDMap map[string]uuid.UUID) error {
	for _, n := range typeNames {
		mtid, ok := typeIDMap[n]
		if !ok {
			return fmt.Errorf("type name %q not found in map", n)
		}
		tag := models.MenuTag{MenuItemID: menuItemID, MenuTypeID: mtid}
		if err := db.Create(&tag).Error; err != nil {
			return fmt.Errorf("create menutag: %w", err)
		}
	}
	return nil
}

func ensureTypes(db *gorm.DB, menuItemID uuid.UUID, typeNames []string, typeIDMap map[string]uuid.UUID) error {
	// ดึง tag ปัจจุบัน
	var tags []models.MenuTag
	if err := db.Where("menu_item_id = ?", menuItemID).Find(&tags).Error; err != nil {
		return fmt.Errorf("list existing tags: %w", err)
	}
	exists := make(map[uuid.UUID]bool, len(tags))
	for _, t := range tags {
		exists[t.MenuTypeID] = true
	}
	// เติมเฉพาะที่ยังไม่มี
	for _, n := range typeNames {
		mtid, ok := typeIDMap[n]
		if !ok {
			return fmt.Errorf("type name %q not found in map", n)
		}
		if exists[mtid] {
			continue
		}
		tag := models.MenuTag{MenuItemID: menuItemID, MenuTypeID: mtid}
		if err := db.Create(&tag).Error; err != nil {
			return fmt.Errorf("create menutag: %w", err)
		}
	}
	return nil
}

func seedTableTimeslots(db *gorm.DB) error {

	// Create tables
	var tables []models.Table
	for col := 1; col <= 3; col++ {
		for row := 1; row <= 6; row++ {
			table := models.Table{
				PeopleNum: 6,
				Row:      fmt.Sprintf("%c", 'A'+(row-1)),
				Col:      fmt.Sprintf("%d", col),
			}
			if err := db.Create(&table).Error; err != nil {
				return err
			}
			tables = append(tables, table)
		}
	}

	// Create time slots from 10:01 to 13:00 with 15-minute intervals
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	baseDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	gap := 1 * time.Minute
	start := baseDate.Add(10 * time.Hour + gap)   // 10:01 in Thailand
	end := baseDate.Add(13 * time.Hour)     // 13:00 in Thailand
	duration := 14 * time.Minute
	var timeSlots []models.Timeslot

	for t := start; t.Before(end); t = t.Add(duration + gap) {
		timeSlot := models.Timeslot{
			StartTime: t,
			EndTime:   t.Add(duration),
		}
		if err := db.Create(&timeSlot).Error; err != nil {
			return err
		}
		timeSlots = append(timeSlots, timeSlot)
	}


	// Create table-time slot associations
	for _, table := range tables {
		for _, timeSlot := range timeSlots {

			status := "available"
			if timeSlot.EndTime.In(loc).After(time.Now().In(loc)) {
				status = "expired"
			}

			tableTimeslot := models.TableTimeslot{
				TableID:    table.ID,
				TimeslotID: timeSlot.ID,
				Status:     status,
			}
			if err := db.Create(&tableTimeslot).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

