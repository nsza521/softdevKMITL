package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	models "backend/internal/db_model"
)

type TableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) *TableRepository {
	return &TableRepository{
		db: db,
	}
}

func (r *TableRepository) CreateTable(table *models.Table) error {
	return r.db.Create(table).Error
}

func (r *TableRepository) GetAllTables() ([]*models.Table, error) {
	var tables []*models.Table
	if err := r.db.Order("created_at ASC").Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *TableRepository) CreateTimeslot(timeslot *models.Timeslot) error {
	return r.db.Create(timeslot).Error
}

func (r *TableRepository) GetAllTimeslots() ([]*models.Timeslot, error) {
	var timeslots []*models.Timeslot
	if err := r.db.Order("created_at ASC").Find(&timeslots).Error; err != nil {
		return nil, err
	}
	return timeslots, nil
}

func (r *TableRepository) GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]*models.TableTimeslot, error) {
	var tableTimeslots []*models.TableTimeslot
	if err := r.db.Where("timeslot_id = ?", timeslotID).Find(&tableTimeslots).Error; err != nil {
		return nil, err
	}
	return tableTimeslots, nil
}