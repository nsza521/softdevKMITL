package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	models "backend/internal/db_model"
	iface "backend/internal/menu/interfaces"
	"gorm.io/gorm/logger"
)

type menuRepo struct{ db *gorm.DB }

func NewMenuRepository(db *gorm.DB) iface.MenuRepository { return &menuRepo{db: db} }

func (r *menuRepo) ListMenuByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuItem, error) {
	var items []models.MenuItem
	err := r.db.WithContext(ctx).
		Model(&models.MenuItem{}).
		Select("menu_items.*").
		Joins("JOIN menu_tags  mt  ON mt.menu_item_id = menu_items.id").
		Joins("JOIN menu_types mty ON mty.id = mt.menu_type_id").
		Where("mty.restaurant_id = ?", restaurantID).
		Group("menu_items.id").
		Preload("MenuTypes", "menu_types.restaurant_id = ?", restaurantID).
		Find(&items).Error
	return items, err
}

func (r *menuRepo) RestaurantExists(ctx context.Context, restaurantID uuid.UUID) error {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Restaurant{}).
		Where("id = ?", restaurantID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

//Get detail 


func (r *menuRepo) GetItemWithTypesAndAddOns(itemID uuid.UUID) (*models.MenuItem, error) {
    var item models.MenuItem
    db := r.db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})

    err := db.
        Preload("MenuTypes").
        Preload("MenuTypes.AddOnGroups").
        Preload("MenuTypes.AddOnGroups.Options").
        Preload("AddOnGroups").
        Preload("AddOnGroups.Options").
        First(&item, "id = ?", itemID).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}




func (r *menuRepo) CreateMenuItem(ctx context.Context, mi *models.MenuItem) error {
	return r.db.WithContext(ctx).Create(mi).Error
}
func (r *menuRepo) UpdateMenuItem(ctx context.Context, id uuid.UUID, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.MenuItem{}).Where("id = ?", id).Updates(fields).Error
}
func (r *menuRepo) DeleteMenuItem(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_item_id = ?", id).Delete(&models.MenuTag{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&models.MenuItem{}).Error
	})
}

func (r *menuRepo) AttachMenuTypes(ctx context.Context, itemID uuid.UUID, typeIDs []uuid.UUID) error {
	if len(typeIDs) == 0 { return nil }
	tags := make([]models.MenuTag, 0, len(typeIDs))
	for _, tid := range typeIDs { tags = append(tags, models.MenuTag{MenuItemID: itemID, MenuTypeID: tid}) }
	return r.db.WithContext(ctx).Create(&tags).Error
}

func (r *menuRepo) ReplaceMenuTypes(ctx context.Context, itemID uuid.UUID, typeIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_item_id = ?", itemID).Delete(&models.MenuTag{}).Error; err != nil {
			return err
		}
		if len(typeIDs) == 0 { return nil }
		tags := make([]models.MenuTag, 0, len(typeIDs))
		for _, tid := range typeIDs { tags = append(tags, models.MenuTag{MenuItemID: itemID, MenuTypeID: tid}) }
		return tx.Create(&tags).Error
	})
}

func (r *menuRepo) VerifyMenuTypesBelongToRestaurant(ctx context.Context, restaurantID uuid.UUID, typeIDs []uuid.UUID) error {
	if len(typeIDs) == 0 { return errors.New("menu_type_ids required") }
	var cnt int64
	if err := r.db.WithContext(ctx).
		Model(&models.MenuType{}).
		Where("restaurant_id = ? AND id IN ?", restaurantID, typeIDs).
		Count(&cnt).Error; err != nil {
		return err
	}
	if cnt != int64(len(typeIDs)) {
		return errors.New("some menu_type_ids do not belong to this restaurant")
	}
	return nil
}

func (r *menuRepo) LoadMenuItemWithTypes(ctx context.Context, id uuid.UUID) (*models.MenuItem, error) {
	var m models.MenuItem
	if err := r.db.WithContext(ctx).
		Preload("MenuTypes").
		First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *menuRepo) GetMenuItemByID(ctx context.Context, id uuid.UUID) (*models.MenuItem, error) {
	var m models.MenuItem
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}