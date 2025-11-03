package seed

import (
	"fmt"
	"os"
	"path/filepath"
	"mime/multipart"
	"net/textproto"

	models "backend/internal/db_model"
	"backend/internal/utils"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

/*
	--- helpers สำหรับรูปเมนู / รูปร้าน ---
	ของเดิมคุณมี uploadSampleRestaurantImage อยู่แล้วในไฟล์ noodle (แม้จะไม่ได้ paste มา)
	ในไฟล์นี้จะสมมติว่ามีฟังก์ชัน uploadSampleRestaurantImage(minioClient, filename) แบบเดียวกัน

	ผมจะก็อป UploadSampleMenuImage จากของเดิมมาใช้ แต่เปลี่ยนชื่อรูปเป็น ChickenRice_*.png
*/

func UploadSampleChickenMenuImage(minioClient *minio.Client, filename string) (string, error) {
	basePath, _ := os.Getwd()
	filePath := filepath.Join(basePath, "internal", "assets", "images", "Menu", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file error: %v", err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Size:     fileInfo.Size(),
		Header:   make(textproto.MIMEHeader),
	}
	fileHeader.Header.Set("Content-Type", "image/png")

	objectName := fmt.Sprintf("menu-items/%s", fileHeader.Filename)

	url, err := utils.UploadImage(file, fileHeader, "restaurant-pictures", objectName, minioClient)
	if err != nil {
		return "", err
	}
	return url, nil
}

// เรียกจาก InitAllSeedData(db) หลัง seed ร้านพื้นฐานแล้ว
func seedFixedForChickenRiceShop(db *gorm.DB, minioClient *minio.Client) error {
	rest, created, err := getOrCreateChickenRiceRestaurant(db)
	if err != nil {
		return err
	}
	if created {
		fmt.Println("created chicken rice shop:", rest.Username)
	}

	// ---- อัปโหลดรูปโปรไฟล์ร้าน ----
	filename := "restaurant_chicken_rice.png"
	if rest.ProfilePic == nil || *rest.ProfilePic == "" {
		imgURL, err := uploadSampleRestaurantImage(minioClient, filename)
		if err != nil {
			return fmt.Errorf("upload chicken rice shop image failed: %w", err)
		}
		rest.ProfilePic = &imgURL
		if err := db.Save(rest).Error; err != nil {
			return fmt.Errorf("save chicken rice shop image URL: %w", err)
		}
		fmt.Println("Uploaded chicken rice shop image:", imgURL)
	}

	// ถ้าในร้านนี้มี tag อยู่แล้ว แปลว่าเคย seed ไปแล้ว → ข้ามเมนูหลัก เหลือแค่เติม Add-on ให้ชัวร์
	var cnt int64
	if err := db.
		Table("menu_tags").
		Joins("JOIN menu_types mty ON mty.id = menu_tags.menu_type_id").
		Where("mty.restaurant_id = ?", rest.ID).
		Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return seedAddOnsForChickenRiceShop(db, rest.ID)
	}

	// เคลียร์ข้อมูลเมนูเก่าเฉพาะร้านนี้ (กันซ้ำ)
	if err := resetMenusForRestaurant(db, rest.ID); err != nil {
		return err
	}

	// ---------- 1) หมวดหมู่เมนูของร้านข้าวมันไก่ ----------
	/*
		หมวดที่จะใช้:
		- "ข้าวมันไก่"
		- "กับข้าว / อื่นๆ"
		- "ทานเล่น"
		- "เครื่องดื่ม"
	*/
	typeNames := []string{
		"ข้าวมันไก่",
		"กับข้าว / อื่นๆ",
	}
	typeID := make(map[string]uuid.UUID, len(typeNames))
	for _, n := range typeNames {
		mt, err := getOrCreateMenuType(db, rest.ID, n)
		if err != nil {
			return err
		}
		typeID[n] = mt.ID
	}

	// ---------- 2) เมนูหลักของร้านข้าวมันไก่ ----------
	type itemSpec struct {
		Name        string
		Price       float64
		TimeTaken   int
		Description string
		Types       []string // ใส่ชื่อหมวดเพื่อ Tag
	}
	items := []itemSpec{
		{
			Name:        "ข้าวมันไก่ต้ม",
			Price:       50,
			TimeTaken:   4,
			Description: "ข้าวมันหอม ไก่ต้มเนื้อนุ่ม เสิร์ฟพร้อมน้ำจิ้มเต้าเจี้ยว",
			Types:       []string{"ข้าวมันไก่"},
		},
		{
			Name:        "ข้าวมันไก่ทอด",
			Price:       55,
			TimeTaken:   5,
			Description: "ไก่ทอดกรอบ ชิ้นหนา น้ำจิ้มซีอิ๊วพริก+กระเทียม",
			Types:       []string{"ข้าวมันไก่"},
		},
		{
			Name:        "ข้าวมันไก่รวม",
			Price:       60,
			TimeTaken:   6,
			Description: "ไก่ต้ม + ไก่ทอด ในจานเดียว",
			Types:       []string{"ข้าวมันไก่"},
		},
		{
			Name:        "ไก่สับจาน",
			Price:       80,
			TimeTaken:   5,
			Description: "ไก่ต้มสับล้วนๆ ไม่เอาข้าว",
			Types:       []string{"กับข้าว / อื่นๆ"},
		},
		{
			Name:        "ต้มซุปฟักกระดูกหมู",
			Price:       30,
			TimeTaken:   3,
			Description: "ซุปร้อนๆ หวานน้ำต้มกระดูกหมู ใส่ฟักนุ่ม",
			Types:       []string{"กับข้าว / อื่นๆ"},
		},
	}

	// ---------- 3) สร้าง MenuItem + รูป + Tag หมวด ----------
	for i, it := range items {

		imgFilename := fmt.Sprintf("MenuChicken_%d.png", i+1)
		imgURL, err := UploadSampleChickenMenuImage(minioClient, imgFilename)
		if err != nil {
			return fmt.Errorf("failed to upload image for %s: %v", it.Name, err)
		}

		mi := models.MenuItem{
			RestaurantID: rest.ID,
			Name:         it.Name,
			Price:        it.Price,
			TimeTaken:    it.TimeTaken,
			Description:  it.Description,
			MenuPic:      &imgURL,
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
				return fmt.Errorf("type %q not found for chicken rice shop", tn)
			}
			if err := db.Create(&models.MenuTag{
				MenuItemID: mi.ID,
				MenuTypeID: mtid,
			}).Error; err != nil {
				return fmt.Errorf("tag %q for %q: %w", tn, it.Name, err)
			}
		}
	}

	// ---------- 4) เติม Add-on ----------
	if err := seedAddOnsForChickenRiceShop(db, rest.ID); err != nil {
		return err
	}

	// ---------- 5) เคลียร์เมนูที่ไม่มี tag (กันขยะ) ----------
	if err := deleteOrphanMenuItems(db); err != nil {
		return err
	}

	return nil
}

// ================== เติม Add-on ของร้านข้าวมันไก่ ==================

func seedAddOnsForChickenRiceShop(db *gorm.DB, restaurantID uuid.UUID) error {
	// Map หมวดที่จะให้ Add-on ใช้งาน
	typeMap, err := getMenuTypeIDsByNames(db, restaurantID, []string{
		"ข้าวมันไก่",
		"กับข้าว / อื่นๆ",
	})
	if err != nil {
		return err
	}

	// 1) กลุ่ม "เลือกน้ำซุป"
	soupGroup, err := ensureAddOnGroup(db, restaurantID, ensureGroupInput{
		Name:      "เลือกน้ำซุป",
		Required:  true,
		MinSelect: intPtr(1),
		MaxSelect: intPtr(1),
		AllowQty:  false,
	})
	if err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, soupGroup.ID, ensureOptionInput{
		Name:       "ซุปฟัก",
		PriceDelta: 0,
		IsDefault:  true,
	}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, soupGroup.ID, ensureOptionInput{
		Name:       "ซุปใส",
		PriceDelta: 0,
	}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, soupGroup.ID, ensureOptionInput{
		Name:       "ไม่เอาซุป",
		PriceDelta: 0,
	}); err != nil {
		return err
	}

	// ผูกกลุ่มน้ำซุปเข้ากับเมนู "ข้าวมันไก่"
	if mtid, ok := typeMap["ข้าวมันไก่"]; ok {
		if err := ensureTypeGroupLink(db, mtid, soupGroup.ID); err != nil {
			return err
		}
	}

	// 2) กลุ่ม "เพิ่มไก่"
	extraChickenGroup, err := ensureAddOnGroup(db, restaurantID, ensureGroupInput{
		Name:      "เพิ่มไก่",
		Required:  false,
		MinSelect: intPtr(0),
		MaxSelect: intPtr(3),
		AllowQty:  true, // ให้เลือกได้หลายชิ้น/หลายส่วน
	})
	if err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, extraChickenGroup.ID, ensureOptionInput{
		Name:       "ไก่ต้มเพิ่ม",
		PriceDelta: 20,
		MaxQty:     intPtr(3),
	}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, extraChickenGroup.ID, ensureOptionInput{
		Name:       "ไก่ทอดเพิ่ม",
		PriceDelta: 25,
		MaxQty:     intPtr(3),
	}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, extraChickenGroup.ID, ensureOptionInput{
		Name:       "หนังไก่ทอด",
		PriceDelta: 15,
		MaxQty:     intPtr(2),
	}); err != nil {
		return err
	}

	// ผูกกลุ่ม "เพิ่มไก่" กับหมวด "ข้าวมันไก่"
	if mtid, ok := typeMap["ข้าวมันไก่"]; ok {
		if err := ensureTypeGroupLink(db, mtid, extraChickenGroup.ID); err != nil {
			return err
		}
	}

	// 3) กลุ่ม "น้ำจิ้ม"
	sauceGroup, err := ensureAddOnGroup(db, restaurantID, ensureGroupInput{
		Name:      "น้ำจิ้ม",
		Required:  true,
		MinSelect: intPtr(1),
		MaxSelect: intPtr(2),
		AllowQty:  false,
	})
	if err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, sauceGroup.ID, ensureOptionInput{
		Name:       "เต้าเจี้ยวขิง",
		PriceDelta: 0,
		IsDefault:  true,
	}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, sauceGroup.ID, ensureOptionInput{
		Name:       "ซีอิ๊วดำพริกกระเทียม",
		PriceDelta: 0,
	}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, sauceGroup.ID, ensureOptionInput{
		Name:       "น้ำจิ้มเผ็ด",
		PriceDelta: 0,
	}); err != nil {
		return err
	}

	// ผูกน้ำจิ้มกับ "ข้าวมันไก่"
	if mtid, ok := typeMap["ข้าวมันไก่"]; ok {
		if err := ensureTypeGroupLink(db, mtid, sauceGroup.ID); err != nil {
			return err
		}
	}

	// 4) กลุ่ม "จานแชร์/กับข้าว"
	// ใช้สำหรับหมวด "กับข้าว / อื่นๆ" เช่น ไก่สับจาน
	shareGroup, err := ensureAddOnGroup(db, restaurantID, ensureGroupInput{
		Name:      "จานแชร์",
		Required:  false,
		MinSelect: intPtr(0),
		MaxSelect: intPtr(1),
		AllowQty:  false,
	})
	if err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, shareGroup.ID, ensureOptionInput{
		Name:       "เพิ่มข้าวมัน (ถ้วยเล็ก)",
		PriceDelta: 10,
	}); err != nil {
		return err
	}
	if _, err := ensureAddOnOption(db, shareGroup.ID, ensureOptionInput{
		Name:       "เพิ่มข้าวมัน (ถ้วยใหญ่)",
		PriceDelta: 15,
	}); err != nil {
		return err
	}

	if mtid, ok := typeMap["กับข้าว / อื่นๆ"]; ok {
		if err := ensureTypeGroupLink(db, mtid, shareGroup.ID); err != nil {
			return err
		}
	}

	// ตัวอย่างปิดสต๊อก:
	// _ = markOptionOutOfStock(db, extraChickenGroup.ID, "หนังไก่ทอด", time.Now().Add(6*time.Hour))

	return nil
}

// ================== Restaurant helper เฉพาะร้านข้าวมันไก่ ==================

func getOrCreateChickenRiceRestaurant(db *gorm.DB) (*models.Restaurant, bool, error) {
	const username = "restaurant_chicken_rice"
	const email = "chickenrice@example.com"
	const name = "ร้านข้าวมันไก่"

	var rest models.Restaurant
	if err := db.Where("username = ? OR email = ? OR name = ?", username, email, name).
		First(&rest).Error; err == nil {
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
