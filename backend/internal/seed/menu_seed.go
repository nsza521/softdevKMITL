package seed

import (
	"fmt"
	"gorm.io/gorm"

	"backend/internal/db_model"
	"github.com/google/uuid"
)

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