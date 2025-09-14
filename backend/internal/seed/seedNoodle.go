// seed_noodle_shop.go
package seed

import (
	"fmt"

	models "backend/internal/db_model"
	"backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// เรียกจาก InitAllSeedData(db) หลัง seed ร้านพื้นฐานแล้ว
func seedFixedForNoodleShop(db *gorm.DB) error {
	rest, created, err := getOrCreateNoodleRestaurant(db)
	if err != nil {
		return err
	}
	if created {
		fmt.Println("created noodle shop:", rest.Username)
	}

	var cnt int64
    if err := db.
        Table("menu_tags").
        Joins("JOIN menu_types mty ON mty.id = menu_tags.menu_type_id").
        Where("mty.restaurant_id = ?", rest.ID).
        Count(&cnt).Error; err != nil {
        return err
    }
    if cnt > 0 {
        return nil // มีของอยู่แล้ว ไม่ยุ่ง
    }

	// ล้างข้อมูลเมนูของร้านนี้ก่อน (กันซ้ำ/ข้อมูลเก่า)
	if err := resetMenusForRestaurant(db, rest.ID); err != nil {
		return err
	}

	// สร้างประเภทเมนูเฉพาะร้านก๋วยเตี๋ยว
	typeNames := []string{
		"เมนูเส้น",
		"ก๋วยเตี๋ยวหมูน้ำตก",
		"ก๋วยเตี๋ยวต้มยำ",
	}
	typeID := make(map[string]uuid.UUID, len(typeNames))
	for _, n := range typeNames {
		mt, err := getOrCreateMenuType(db, rest.ID, n)
		if err != nil {
			return err
		}
		typeID[n] = mt.ID
	}

	// เมนูของร้านก๋วยเตี๋ยว (กำหนดตายตัว)
	type itemSpec struct {
		Name        string
		Price       float64
		TimeTaken   int
		Description string
		Types       []string // จะ tag ให้กับ MenuType ของร้านนี้เท่านั้น
	}
	items := []itemSpec{
		{"ก๋วยเตี๋ยวน้ำตกหมู", 50, 5, "เส้นเล็ก น้ำซุปรสเข้ม", []string{"เมนูเส้น", "ก๋วยเตี๋ยวหมูน้ำตก"}},
		{"ก๋วยเตี๋ยวต้มยำกุ้ง", 70, 6, "เผ็ดเปรี้ยวกลมกล่อม", []string{"เมนูเส้น", "ก๋วยเตี๋ยวต้มยำ"}},
		{"ก๋วยเตี๋ยวเนื้อตุ๋น", 80, 7, "เนื้อตุ๋นเปื่อยนุ่ม", []string{"เมนูเส้น"}},
		{"เกี๊ยวทอด", 35, 3, "กรอบอร่อย", []string{}},
		{"ชาดำเย็น", 25, 2, "หวานหอม ดับกระหาย", []string{}},
	}

	// สร้าง MenuItem ใหม่ “สำหรับร้านนี้” และผูก Tag เฉพาะ type ของร้านนี้
	for _, it := range items {
		mi := models.MenuItem{
			Name:        it.Name,
			Price:       it.Price,
			TimeTaken:   it.TimeTaken,
			Description: it.Description,
		}
		if mi.TimeTaken == 0 {
			mi.TimeTaken = 1
		}
		if err := db.Create(&mi).Error; err != nil {
			return fmt.Errorf("create menuitem %q: %w", it.Name, err)
		}
		for _, tn := range it.Types {
			mtid, ok := typeID[tn]
			if !ok {
				return fmt.Errorf("type %q not found for noodle shop", tn)
			}
			if err := db.Create(&models.MenuTag{
				MenuItemID: mi.ID,
				MenuTypeID: mtid,
			}).Error; err != nil {
				return fmt.Errorf("tag %q for %q: %w", tn, it.Name, err)
			}
		}
	}

	// เคลียร์เมนูที่ไม่มี tag ค้าง (ถ้ามี)
	if err := deleteOrphanMenuItems(db); err != nil {
		return err
	}

	return nil
}

// ---------- helpers ----------

// สร้าง/หา restaurant “ร้านก๋วยเตี๋ยว” (username กำหนดเอง)
func getOrCreateNoodleRestaurant(db *gorm.DB) (*models.Restaurant, bool, error) {
	const username = "restaurant_noodle"
	const email = "noodle@example.com"

	var rest models.Restaurant
	if err := db.Where("username = ? OR email = ?", username, email).First(&rest).Error; err == nil {
		return &rest, false, nil
	}

	hashed, err := utils.HashPassword("SecureP@ssw0rd")
	if err != nil {
		return nil, false, err
	}
	rest = models.Restaurant{
		Username: username,
		Email:    email,
		Password: hashed,
	}
	if err := db.Create(&rest).Error; err != nil {
		return nil, false, err
	}
	return &rest, true, nil
}

// ล้างข้อมูลเมนู “เฉพาะร้านนี้”: ลบ tag ของ type ในร้านนี้ + ลบ type ของร้านนี้
func resetMenusForRestaurant(db *gorm.DB, restaurantID uuid.UUID) error {
	var typeIDs []uuid.UUID
	if err := db.Model(&models.MenuType{}).
		Where("restaurant_id = ?", restaurantID).
		Pluck("id", &typeIDs).Error; err != nil {
		return err
	}
	if len(typeIDs) == 0 {
		return nil
	}
	if err := db.Where("menu_type_id IN ?", typeIDs).Delete(&models.MenuTag{}).Error; err != nil {
		return err
	}
	if err := db.Where("id IN ?", typeIDs).Delete(&models.MenuType{}).Error; err != nil {
		return err
	}
	return nil
}

// ลบ menu_items ที่ไม่มี tag เหลือ (กันขยะ)
func deleteOrphanMenuItems(db *gorm.DB) error {
	var orphanIDs []uuid.UUID
	if err := db.Raw(`
		SELECT mi.id
		FROM menu_items mi
		LEFT JOIN menu_tags mt ON mt.menu_item_id = mi.id
		WHERE mt.menu_item_id IS NULL
	`).Scan(&orphanIDs).Error; err != nil {
		return err
	}
	if len(orphanIDs) == 0 {
		return nil
	}
	return db.Where("id IN ?", orphanIDs).Delete(&models.MenuItem{}).Error
}
