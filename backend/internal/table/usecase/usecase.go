package usecase

import (
	"fmt"

	"github.com/google/uuid"

	"backend/internal/db_model"
	"backend/internal/table/dto"
	"backend/internal/table/interfaces"
)

type TableUsecase struct {
	tableRepository interfaces.TableRepository
}

func NewTableUsecase(tableRepository interfaces.TableRepository) interfaces.TableUsecase {
	return &TableUsecase{
		tableRepository: tableRepository,
	}
}

// Table Usecase
func (u *TableUsecase) CreateTable(table *models.Table) error {
	return u.tableRepository.CreateTable(table)
}

func (u *TableUsecase) GetAllTables() ([]dto.TableDetail, error) {
	tables, err := u.tableRepository.GetAllTables()
	if err != nil {
		return nil, err
	}

	var details []dto.TableDetail
	for _, t := range tables {
		detail := dto.TableDetail{
			// ID:       t.ID,
			Row:      t.Row,
			Col:      t.Col,
			MaxSeats: t.MaxSeats,
		}
		details = append(details, detail)
	}

	return details, nil
}

// Timeslot Usecase
func (u *TableUsecase) CreateTimeslot(timeslot *models.Timeslot) error {
	return u.tableRepository.CreateTimeslot(timeslot)
}

func (u *TableUsecase) GetAllTimeslots() ([]dto.TimeslotDetail, error) {
	timeslots, err := u.tableRepository.GetAllTimeslots()
	if err != nil {
		return nil, err
	}

	var details []dto.TimeslotDetail
	for _, t := range timeslots {
		detail := dto.TimeslotDetail{
			ID:        t.ID,
			StartTime: t.StartTime.Format("15:04"),
			EndTime:   t.EndTime.Format("15:04"),
		}
		details = append(details, detail)
	}

	return details, nil
}

// TableTimeslot Usecase
func (u *TableUsecase) GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]dto.TableTimeslotDetail, error) {

	tableTimeslots, err := u.tableRepository.GetTableTimeslotByTimeslotID(timeslotID)
	if err != nil {
		return nil, err
	}

	var details []dto.TableTimeslotDetail

	for i, t := range tableTimeslots {

		table, err := u.tableRepository.GetTableByID(t.TableID)
		if err != nil {
			return nil, fmt.Errorf("failed to get table at index %d: %v", i, err)
		}
		tableDetail := dto.TableDetail{
			// ID:  table.ID,
			Row: table.Row,
			Col: table.Col,
			MaxSeats: table.MaxSeats,
		}

		// timeslot, err := u.tableRepository.GetTimeslotByID(t.TimeslotID)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to get timeslot at index %d: %v", i, err)
		// }
		// timeslotDetail := dto.TimeslotDetail{
		// 	ID:        timeslot.ID,
		// 	StartTime: timeslot.StartTime.Format("15:04"),
		// 	EndTime:   timeslot.EndTime.Format("15:04"),
		// }

		detail := dto.TableTimeslotDetail{
			ID:             t.ID,
			Table:      	tableDetail,
			// TimeslotID:   	t.TimeslotID,
			// Timeslot:   	timeslotDetail,
			Status:       	t.Status,
			ReservedSeats:  t.ReservedSeats,
			// MaxSeats:    	table.MaxSeats,
		}
		details = append(details, detail)
	}

	return details, nil
}

func (u *TableUsecase) GetTableTimeslotByID(id uuid.UUID) (*dto.TableTimeslotDetail, error) {
	tableTimeslot, err := u.tableRepository.GetTableTimeslotByID(id)
	if err != nil {
		return nil, err
	}

	table, err := u.tableRepository.GetTableByID(tableTimeslot.TableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %v", err)
	}
	tableDetail := dto.TableDetail{
		// ID:  table.ID,
		TableRow: table.TableRow,
		TableCol: table.TableCol,
		MaxSeats: table.MaxSeats,
	}

	detail := &dto.TableTimeslotDetail{
		ID:             tableTimeslot.ID,
		Table:      	tableDetail,
		Status:       	tableTimeslot.Status,
		ReservedSeats:  tableTimeslot.ReservedSeats,
		// MaxSeats:    	table.MaxSeats,
	}

	return detail, nil
}