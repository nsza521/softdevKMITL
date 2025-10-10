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

// Table Repository
func (r *TableRepository) CreateTable(table *models.Table) error {
	return r.db.Create(table).Error
}

func (r *TableRepository) GetAllTables() ([]models.Table, error) {
	var tables []models.Table
	if err := r.db.Order("created_at ASC").Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *TableRepository) GetTableByID(id uuid.UUID) (*models.Table, error) {
	var table models.Table
	if err := r.db.First(&table, id).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *TableRepository) UpdateTable(table *models.Table) error {
	if err := r.db.Save(table).Error; err != nil {
		return err
	}
	return nil
}

// Timeslot Repository
func (r *TableRepository) CreateTimeslot(timeslot *models.Timeslot) error {
	return r.db.Create(timeslot).Error
}

func (r *TableRepository) GetAllTimeslots() ([]models.Timeslot, error) {
	var timeslots []models.Timeslot
	if err := r.db.Order("created_at ASC").Find(&timeslots).Error; err != nil {
		return nil, err
	}
	return timeslots, nil
}

func (r *TableRepository) GetTimeslotByID(id uuid.UUID) (*models.Timeslot, error) {
	var timeslot models.Timeslot
	if err := r.db.First(&timeslot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &timeslot, nil
}

func (r *TableRepository) UpdateTimeslot(timeslot *models.Timeslot) error {
	if err := r.db.Save(timeslot).Error; err != nil {
		return err
	}
	return nil
}

func (r *TableRepository) IsTimeslotAvailable(id uuid.UUID) (bool, error) {
	
	tableTimeslots, err := r.GetTableTimeslotByTimeslotID(id)
	if err != nil {
		return false, err
	}

	for _, t := range tableTimeslots {
		if t.Status == "available" || t.Status == "partial" {
			return true, nil
		}
	}
	return false, nil
}

// TableTimeslot Repository
func (r *TableRepository) GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]models.TableTimeslot, error) {
	var tableTimeslots []models.TableTimeslot
	if err := r.db.Where("timeslot_id = ?", timeslotID).Order("created_at ASC").Find(&tableTimeslots).Error; err != nil {
		return nil, err
	}
	return tableTimeslots, nil
}

func (r *TableRepository) GetTableTimeslotByID(id uuid.UUID) (*models.TableTimeslot, error) {
	var tableTimeslot models.TableTimeslot
	if err := r.db.First(&tableTimeslot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tableTimeslot, nil
}

func (r *TableRepository) UpdateTableTimeslot(tableTimeslot *models.TableTimeslot) error {
	if err := r.db.Save(tableTimeslot).Error; err != nil {
		return err
	}
	return nil
}