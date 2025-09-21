package repository

import (
	"backend/internal/menu/interfaces"
	models "backend/internal/db_model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type addOnRepository struct {
	db *gorm.DB
}

func NewAddOnRepository(db *gorm.DB) interfaces.AddOnRepository {
	return &addOnRepository{db: db}
}

// Groups
func (r *addOnRepository) CreateGroup(group *models.MenuAddOnGroup) error {
	return r.db.Create(group).Error
}

func (r *addOnRepository) GetGroupsByRestaurant(restaurantID uuid.UUID) ([]models.MenuAddOnGroup, error) {
	var groups []models.MenuAddOnGroup
	err := r.db.Preload("Options").Where("restaurant_id = ?", restaurantID).Find(&groups).Error
	return groups, err
}

func (r *addOnRepository) GetGroupByID(id uuid.UUID) (*models.MenuAddOnGroup, error) {
	var group models.MenuAddOnGroup
	err := r.db.Preload("Options").First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *addOnRepository) UpdateGroup(group *models.MenuAddOnGroup) error {
	return r.db.Save(group).Error
}

func (r *addOnRepository) DeleteGroup(id uuid.UUID) error {
	return r.db.Delete(&models.MenuAddOnGroup{}, "id = ?", id).Error
}

// Options
func (r *addOnRepository) CreateOption(opt *models.MenuAddOnOption) error {
	return r.db.Create(opt).Error
}

func (r *addOnRepository) GetOptionByID(id uuid.UUID) (*models.MenuAddOnOption, error) {
	var opt models.MenuAddOnOption
	err := r.db.First(&opt, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

func (r *addOnRepository) UpdateOption(opt *models.MenuAddOnOption) error {
	return r.db.Save(opt).Error
}

func (r *addOnRepository) DeleteOption(id uuid.UUID) error {
	return r.db.Delete(&models.MenuAddOnOption{}, "id = ?", id).Error
}
