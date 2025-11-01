// seed_noodle_shop.go
package seed

import (
	"fmt"
	"time"

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

	// ถ้ามี tag ของร้านนี้อยู่แล้ว แปลว่าน่าจะ seed ไปแล้ว — ข้าม
	var cnt int64
	if err := db.
		Table("menu_tags").
		Joins("JOIN menu_types mty ON mty.id = menu_tags.menu_type_id").
		Where("mty.restaurant_id = ?", rest.ID).
		Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		// แต่เผื่อยังไม่มี Add-on ให้เติมเฉพาะส่วน Add-on ได้ (idempotent)
		return seedAddOnsForNoodleShop(db, rest.ID)
	}

	// ล้างข้อมูลเมนูของร้านนี้ก่อน (กันซ้ำ/ข้อมูลเก่า)
	if err := resetMenusForRestaurant(db, rest.ID); err != nil {
		return err
	}

	// 1) สร้างประเภทเมนูเฉพาะร้านก๋วยเตี๋ยว
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

	// 2) เมนูของร้านก๋วยเตี๋ยว (กำหนดตายตัว)
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

	// 3) สร้าง MenuItem ใหม่ “สำหรับร้านนี้” และผูก Tag เฉพาะ type ของร้านนี้
	for _, it := range items {
		mi := models.MenuItem{
			RestaurantID: rest.ID, // << สำคัญ: ลงร้านนี้
			Name:         it.Name,
			Price:        it.Price,
			TimeTaken:    it.TimeTaken,
			Description:  it.Description,
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

	// 4) เติม Add-on ให้ครบ
	if err := seedAddOnsForNoodleShop(db, rest.ID); err != nil {
		return err
	}

	// เคลียร์เมนูที่ไม่มี tag ค้าง (ถ้ามี)
	if err := deleteOrphanMenuItems(db); err != nil {
		return err
	}

	return nil
}

// เติม Add-on แบบ idempotent สำหรับร้านก๋วยเตี๋ยว
func seedAddOnsForNoodleShop(db *gorm.DB, restaurantID uuid.UUID) error {
	// เตรียมอ้างอิงประเภทที่ต้องผูก Add-on
	typeMap, err := getMenuTypeIDsByNames(db, restaurantID, []string{
		"เมนูเส้น",
		"ก๋วยเตี๋ยวหมูน้ำตก",
		"ก๋วยเตี๋ยวต้มยำ",
	})
	if err != nil {
		return err
	}

	// === 1) กลุ่ม "เลือกเส้น" ===
	noodleGroup, err := ensureAddOnGroup(db, restaurantID, ensureGroupInput{
		Name:      "เลือกเส้น",
		Required:  true,
		MinSelect: intPtr(1),
		MaxSelect: intPtr(1),
		AllowQty:  false,
	})
	if err != nil {
		return err
	}
	// ตัวเลือกใน “เลือกเส้น”
	if _, err := ensureAddOnOption(db, noodleGroup.ID, ensureOptionInput{Name: "เส้นเล็ก", PriceDelta: 0, IsDefault: true}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, noodleGroup.ID, ensureOptionInput{Name: "เส้นใหญ่", PriceDelta: 0}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, noodleGroup.ID, ensureOptionInput{Name: "บะหมี่", PriceDelta: 5}); err != nil {
		return err
	}
	// ผูกเข้ากับหมวดที่เกี่ยวข้อง
	for _, tn := range []string{"เมนูเส้น", "ก๋วยเตี๋ยวหมูน้ำตก", "ก๋วยเตี๋ยวต้มยำ"} {
		if mtid, ok := typeMap[tn]; ok {
			if err := ensureTypeGroupLink(db, mtid, noodleGroup.ID); err != nil {
				return err
			}
		}
	}

	// === 2) กลุ่ม "ระดับความเผ็ด" ===
	spiceGroup, err := ensureAddOnGroup(db, restaurantID, ensureGroupInput{
		Name:      "ระดับความเผ็ด",
		Required:  true,
		MinSelect: intPtr(1),
		MaxSelect: intPtr(1),
		AllowQty:  false,
	})
	if err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, spiceGroup.ID, ensureOptionInput{Name: "ไม่เผ็ด", PriceDelta: 0, IsDefault: true}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, spiceGroup.ID, ensureOptionInput{Name: "เผ็ดกลาง", PriceDelta: 0}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, spiceGroup.ID, ensureOptionInput{Name: "เผ็ดมาก", PriceDelta: 0}); err != nil {
		return err
	}
	for _, tn := range []string{"เมนูเส้น", "ก๋วยเตี๋ยวหมูน้ำตก", "ก๋วยเตี๋ยวต้มยำ"} {
		if mtid, ok := typeMap[tn]; ok {
			if err := ensureTypeGroupLink(db, mtid, spiceGroup.ID); err != nil {
				return err
			}
		}
	}

	// === 3) กลุ่ม "ท็อปปิ้ง/เพิ่ม" (เลือกได้หลายอย่าง + ระบุจำนวนได้) ===
	toppingGroup, err := ensureAddOnGroup(db, restaurantID, ensureGroupInput{
		Name:      "ท็อปปิ้งเพิ่ม",
		Required:  false,
		MinSelect: intPtr(0),
		MaxSelect: intPtr(5),
		AllowQty:  true, // ใส่ +1 +2 …
	})
	if err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, toppingGroup.ID, ensureOptionInput{Name: "ไข่ลวก", PriceDelta: 10, MaxQty: intPtr(3)}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, toppingGroup.ID, ensureOptionInput{Name: "ลูกชิ้นหมูเพิ่ม", PriceDelta: 15, MaxQty: intPtr(5)}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, toppingGroup.ID, ensureOptionInput{Name: "หมูเพิ่ม", PriceDelta: 20, MaxQty: intPtr(3)}); err != nil {
		return err
	}
	// ผูกเฉพาะหมวดเส้น/ก๋วยเตี๋ยว
	for _, tn := range []string{"เมนูเส้น", "ก๋วยเตี๋ยวหมูน้ำตก", "ก๋วยเตี๋ยวต้มยำ"} {
		if mtid, ok := typeMap[tn]; ok {
			if err := ensureTypeGroupLink(db, mtid, toppingGroup.ID); err != nil {
				return err
			}
		}
	}

	// ตัวอย่างปิดสต๊อกบาง option ชั่วคราว (ถ้า model มี IsAvailable/OutOfStockUntil)
	// _ = markOptionOutOfStock(db, noodleGroup.ID, "บะหมี่", time.Now().Add(12*time.Hour))

	return nil
}


// ---------- helpers ----------

// สร้าง/หา restaurant “ร้านก๋วยเตี๋ยว” (username กำหนดเอง)
func getOrCreateNoodleRestaurant(db *gorm.DB) (*models.Restaurant, bool, error) {
	const username = "restaurant_noodle"
	const email = "noodle@example.com"
	const name = "ร้านก๋วยเตี๋ยว"

	var rest models.Restaurant
	if err := db.Where("username = ? OR email = ? OR name = ?", username, email, name).First(&rest).Error; err == nil {
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
		Name:     name,
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

// ==== Add-on helpers ====

// ป้องกัน duplicate กลุ่ม
type ensureGroupInput struct {
	Name      string
	Required  bool
	MinSelect *int
	MaxSelect *int
	AllowQty  bool
}

func ensureAddOnGroup(db *gorm.DB, restaurantID uuid.UUID, in ensureGroupInput) (*models.MenuAddOnGroup, error) {
	var g models.MenuAddOnGroup
	err := db.Where("restaurant_id = ? AND name = ?", restaurantID, in.Name).First(&g).Error
	if err == nil {
		// อัพเดตค่าเผื่อ schema เปลี่ยน
		changed := false
		if g.Required != in.Required { g.Required = in.Required; changed = true }
		if (g.MinSelect == nil && in.MinSelect != nil) || (g.MinSelect != nil && in.MinSelect == nil) || (g.MinSelect != nil && in.MinSelect != nil && *g.MinSelect != *in.MinSelect) {
			g.MinSelect = in.MinSelect; changed = true
		}
		if (g.MaxSelect == nil && in.MaxSelect != nil) || (g.MaxSelect != nil && in.MaxSelect == nil) || (g.MaxSelect != nil && in.MaxSelect != nil && *g.MaxSelect != *in.MaxSelect) {
			g.MaxSelect = in.MaxSelect; changed = true
		}
		if g.AllowQty != in.AllowQty { g.AllowQty = in.AllowQty; changed = true }
		if changed {
			if err := db.Save(&g).Error; err != nil { return nil, err }
		}
		return &g, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	g = models.MenuAddOnGroup{
		RestaurantID: restaurantID,
		Name:         in.Name,
		Required:     in.Required,
		MinSelect:    in.MinSelect,
		MaxSelect:    in.MaxSelect,
		AllowQty:     in.AllowQty,
	}
	if err := db.Create(&g).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

type ensureOptionInput struct {
	Name       string
	PriceDelta float64
	IsDefault  bool
	MaxQty     *int
}

func ensureAddOnOption(db *gorm.DB, groupID uuid.UUID, in ensureOptionInput) (*models.MenuAddOnOption, error) {
	var o models.MenuAddOnOption
	err := db.Where("group_id = ? AND name = ?", groupID, in.Name).First(&o).Error
	if err == nil {
		changed := false
		if o.PriceDelta != in.PriceDelta { o.PriceDelta = in.PriceDelta; changed = true }
		if o.IsDefault != in.IsDefault { o.IsDefault = in.IsDefault; changed = true }
		// compare pointer
		if (o.MaxQty == nil && in.MaxQty != nil) || (o.MaxQty != nil && in.MaxQty == nil) || (o.MaxQty != nil && in.MaxQty != nil && *o.MaxQty != *in.MaxQty) {
			o.MaxQty = in.MaxQty; changed = true
		}
		if changed {
			if err := db.Save(&o).Error; err != nil { return nil, err }
		}
		return &o, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	o = models.MenuAddOnOption{
		GroupID:    groupID,
		Name:       in.Name,
		PriceDelta: in.PriceDelta,
		IsDefault:  in.IsDefault,
		MaxQty:     in.MaxQty,
	}
	if err := db.Create(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// ผูก MenuType ↔ AddOnGroup (many2many: menu_type_addon_groups)
func ensureTypeGroupLink(db *gorm.DB, menuTypeID, groupID uuid.UUID) error {
	// ตรวจซ้ำ
	var n int64
	if err := db.
		Model(&models.MenuTypeAddOnGroup{}).
		Where("menu_type_id = ? AND add_on_group_id = ?", menuTypeID, groupID).
		Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	link := models.MenuTypeAddOnGroup{
		MenuTypeID:   menuTypeID,
		AddOnGroupID: groupID,
	}
	return db.Create(&link).Error
}

// ดึง ids ของ MenuType ตามชื่อ (ภายในร้านเดียวกัน)
func getMenuTypeIDsByNames(db *gorm.DB, restaurantID uuid.UUID, names []string) (map[string]uuid.UUID, error) {
	type rec struct {
		ID   uuid.UUID
		Type string
	}
	var rows []rec
	if err := db.
		Model(&models.MenuType{}).
		Where("restaurant_id = ? AND type IN ?", restaurantID, names).
		Select("id, type").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	res := make(map[string]uuid.UUID, len(rows))
	for _, r := range rows {
		res[r.Type] = r.ID
	}
	return res, nil
}

// (ตัวอย่าง) ปิดสต๊อก option ชั่วคราว ถ้าโมเดลคุณมีฟิลด์ IsAvailable/OutOfStockUntil
func markOptionOutOfStock(db *gorm.DB, groupID uuid.UUID, optionName string, until time.Time) error {
	// ต้องมีฟิลด์ใน models.MenuAddOnOption ก่อน (เช่น IsAvailable, OutOfStockUntil)
	// ด้านล่างเป็นตัวอย่างโค้ด — คอมเมนต์ไว้เผื่อคุณใส่ฟิลด์แล้วค่อยเปิดใช้
	/*
		return db.Model(&models.MenuAddOnOption{}).
			Where("group_id = ? AND name = ?", groupID, optionName).
			Updates(map[string]interface{}{
				"is_available":       false,
				"out_of_stock_until": until,
			}).Error
	*/
	return nil
}

// ==== utils ====
func intPtr(v int) *int { return &v }