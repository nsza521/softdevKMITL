package repository

import (
	"gorm.io/gorm"

	"github.com/google/uuid"

	models "backend/internal/db_model"
	// restint "backend/internal/restaurant/interfaces"
)

type RestaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{
		db: db,
	}
}

func (r *RestaurantRepository) Create(restaurant *models.Restaurant) (*models.Restaurant, error) {

	result := r.db.Create(restaurant)
	if result.Error != nil {
		return nil, result.Error
	}

	restaurant, err := r.GetByUsername(restaurant.Username)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (r *RestaurantRepository) IsRestaurantExists(username string, email string) (bool, error) {

	var restaurant models.Restaurant

	// Check if name exists
	result := r.db.First(&restaurant, "username = ?", username)
	if result.Error == nil {
		return true, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}

	// Check if email exists
	result = r.db.First(&restaurant, "email = ?", email)
	if result.Error == nil {
		return true, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}

	return false, nil
}

func (r *RestaurantRepository) GetByUsername(username string) (*models.Restaurant, error) {

	var restaurant models.Restaurant

	result := r.db.First(&restaurant, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}

func (r *RestaurantRepository) GetByID(id uuid.UUID) (*models.Restaurant, error) {

	var restaurant models.Restaurant

	result := r.db.First(&restaurant, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}

func (r *RestaurantRepository) GetAll() ([]*models.Restaurant, error) {

	var restaurants []*models.Restaurant

	result := r.db.Order("created_at ASC").Find(&restaurants)
	if result.Error != nil {
		return nil, result.Error
	}
	return restaurants, nil
}

func (r *RestaurantRepository) CreateBankAccount(bankAccount *models.BankAccount) error {
	return r.db.Create(bankAccount).Error
}

func (r *RestaurantRepository) Update(restaurant *models.Restaurant) error {
	return r.db.Save(&restaurant).Error
}

// ✅ อัปเดตเฉพาะฟิลด์ที่ != nil
func (r *RestaurantRepository) PartialUpdate(
	restaurantID uuid.UUID,
	name string,
	menuType *string,
	// addOnMenuItems []string,
) (*models.Restaurant, error) {

	// 1) อัปเดต name/menuType เฉพาะที่ != nil
	updates := map[string]any{}
	if name != "" {
		updates["name"] = name
	}
	if menuType != nil {
		updates["menu_type"] = *menuType
	}
	if len(updates) > 0 {
		if err := r.db.Model(&models.Restaurant{}).
			Where("id = ?", restaurantID).
			Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 2) ถ้า addOnMenuItems != nil ให้ “แทนที่ทั้งหมด”
	// if addOnMenuItems != nil {
	// 	// ลบของเก่า
	// 	if err := r.db.Where("restaurant_id = ?", restaurantID).
	// 		Delete(&models.RestaurantAddOn{}).Error; err != nil {
	// 		return nil, err
	// 	}
	// 	// ใส่ของใหม่ (ถ้า [] ว่าง = ล้างหมด)
	// 	if len(addOnMenuItems) > 0 {
	// 		bulk := make([]models.RestaurantAddOn, 0, len(addOnMenuItems))
	// 		for _, it := range addOnMenuItems {
	// 			bulk = append(bulk, models.RestaurantAddOn{
	// 				ID:           uuid.New(),
	// 				RestaurantID: restaurantID,
	// 				ItemName:     it,
	// 			})
	// 		}
	// 		if err := r.db.Create(&bulk).Error; err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	// 3) คืนค่า restaurant ล่าสุด (ให้ตรงกับลายเซ็น 2 ค่า)
	var updated models.Restaurant
	if err := r.db.Where("id = ?", restaurantID).First(&updated).Error; err != nil {
		return nil, err
	}
	return &updated, nil
}

// ✅ ลบของเก่าทั้งหมด แล้วเพิ่มใหม่ (ถ้า [] ว่าง = ล้าง)
func (r *RestaurantRepository) ReplaceAddOnItems(restaurantID uuid.UUID, items []string) error {
	if err := r.db.Where("restaurant_id = ?", restaurantID).Delete(&models.RestaurantAddOn{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	bulk := make([]models.RestaurantAddOn, 0, len(items))
	for _, it := range items {
		bulk = append(bulk, models.RestaurantAddOn{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			ItemName:     it,
		})
	}
	return r.db.Create(&bulk).Error
}

// ดึง Add-on เป็น []string
func (r *RestaurantRepository) GetAddOnMenuItems(restaurantID uuid.UUID) ([]string, error) {
	var rows []models.RestaurantAddOn
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Order("item_name asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, x := range rows {
		out = append(out, x.ItemName)
	}
	return out, nil
}

// แทนที่ Add-on ทั้งหมด (ถ้า [] ว่าง = ล้าง)
func (r *RestaurantRepository) ReplaceAddOnMenuItems(restaurantID uuid.UUID, items []string) error {
	if err := r.db.Where("restaurant_id = ?", restaurantID).
		Delete(&models.RestaurantAddOn{}).Error; err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	bulk := make([]models.RestaurantAddOn, 0, len(items))
	for _, it := range items {
		bulk = append(bulk, models.RestaurantAddOn{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			ItemName:     it,
		})
	}
	return r.db.Create(&bulk).Error
}

func (r *RestaurantRepository) UpdateName(restaurantID uuid.UUID, name string) (*models.Restaurant, error) {
    if err := r.db.Model(&models.Restaurant{}).
        Where("id = ?", restaurantID).
        Update("name", name).Error; err != nil {
        return nil, err
    }

    var updated models.Restaurant
    if err := r.db.First(&updated, "id = ?", restaurantID).Error; err != nil {
        return nil, err
    }

    return &updated, nil
}
