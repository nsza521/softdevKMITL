package usecase

import (
	models "backend/internal/db_model"
	"backend/internal/table/interfaces"

	"github.com/google/uuid"
)

type TableUsecase struct {
	tableRepository interfaces.TableRepository
}

func NewTableUsecase(tableRepository interfaces.TableRepository) interfaces.TableUsecase {
	return &TableUsecase{
		tableRepository: tableRepository,
	}
}

func (u *TableUsecase) CreateTable(table *models.Table) error {
	return u.tableRepository.CreateTable(table)
}

func (u *TableUsecase) GetAllTables() ([]*models.Table, error) {
	return u.tableRepository.GetAllTables()
}

func (u *TableUsecase) CreateTimeslot(timeslot *models.Timeslot) error {
	return u.tableRepository.CreateTimeslot(timeslot)
}

func (u *TableUsecase) GetAllTimeslots() ([]*models.Timeslot, error) {
	return u.tableRepository.GetAllTimeslots()
}

func (u *TableUsecase) GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]*models.TableTimeslot, error) {
	return u.tableRepository.GetTableTimeslotByTimeslotID(timeslotID)
}

// func (u *TableUsecase) GetTableByID(id uint) (*models.Table, error) {
// 	return u.tableRepository.GetTableByID(id)
// }