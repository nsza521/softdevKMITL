package usecase

import (
	"fmt"
	"time"

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


func (u *TableUsecase) isDuplicateRowCol(row, col string, excludeID uuid.UUID) (bool, error) {
	tables, err := u.tableRepository.GetAllTables()
	if err != nil {
		return false, err
	}

	for _, t := range tables {
		if excludeID != uuid.Nil && t.ID == excludeID {
			continue
		}

		if t.TableRow == row && t.TableCol == col {
			return true, nil
		}
	}
	return false, nil
}

func (u *TableUsecase) parseTodayTime(input string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)

	parsed, err := time.ParseInLocation("15:04", input, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format (%s): %v", input, err)
	}

	return time.Date(
		now.Year(), now.Month(), now.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0, loc,
	), nil
}

func (u *TableUsecase)isTimeslotOverlapping(newStart, newEnd time.Time, excludeID uuid.UUID) (bool, error) {
	existing, err := u.tableRepository.GetAllTimeslots()
	if err != nil {
		return false, err
	}

	for _, t := range existing {
		// Cross check to exclude the timeslot with excludeID 
		if excludeID != uuid.Nil && t.ID == excludeID {
			continue
		}

		if newStart.Before(t.EndTime) && newEnd.After(t.StartTime) {
			return true, nil
		}
	}
	return false, nil
}

// Table Usecase
func (u *TableUsecase) CreateTable(request dto.CreateTableRequest) (*dto.TableDetail, error) {

	// validate duplicate row+col
	if dup, err := u.isDuplicateRowCol(request.TableRow, request.TableCol, uuid.Nil); err != nil {
		return nil, fmt.Errorf("failed to check duplicates: %v", err)
	} else if dup {
		return nil, fmt.Errorf("table with row %s and col %s already exists", request.TableRow, request.TableCol)
	}

	table := &models.Table{
		TableRow: request.TableRow,
		TableCol: request.TableCol,
		MaxSeats: request.MaxSeats,
	}

	if err := u.tableRepository.CreateTable(table); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	table, err := u.tableRepository.GetTableByID(table.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created table: %v", err)
	}

	return &dto.TableDetail{
		ID:       table.ID,
		TableRow:      table.TableRow,
		TableCol:      table.TableCol,
		MaxSeats: table.MaxSeats,
	}, nil
}

func (u *TableUsecase) GetAllTables() ([]dto.TableDetail, error) {
	tables, err := u.tableRepository.GetAllTables()
	if err != nil {
		return nil, err
	}

	var details []dto.TableDetail
	for _, t := range tables {
		detail := dto.TableDetail{
			ID:       t.ID,
			TableRow: t.TableRow,
			TableCol: t.TableCol,
			MaxSeats: t.MaxSeats,
		}
		details = append(details, detail)
	}

	return details, nil
}

func (u *TableUsecase) GetTableByID(id uuid.UUID) (*dto.TableDetail, error) {
	table, err := u.tableRepository.GetTableByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %v", err)
	}

	detail := &dto.TableDetail{
		ID:       table.ID,
		TableRow:      table.TableRow,
		TableCol:      table.TableCol,
		MaxSeats: table.MaxSeats,
	}

	return detail, nil
}

func (u *TableUsecase) EditTableDetails(id uuid.UUID, request *dto.EditTableRequest) error {
	table, err := u.tableRepository.GetTableByID(id)
	if err != nil {
		return fmt.Errorf("table not found: %v", err)
	}

	if request.TableRow != "" {
		table.TableRow = request.TableRow
	}
	if request.TableCol != "" {
		table.TableCol = request.TableCol
	}
	if request.MaxSeats != 0 {
		table.MaxSeats = request.MaxSeats
	}
	// validate duplicate row+col
	if dup, err := u.isDuplicateRowCol(table.TableRow, table.TableCol, table.ID); err != nil {
		return fmt.Errorf("failed to check duplicates: %v", err)
	} else if dup {
		return fmt.Errorf("table with row %s and col %s already exists", table.TableRow, table.TableCol)
	}

	if err := u.tableRepository.UpdateTable(table); err != nil {
		return fmt.Errorf("failed to update table: %v", err)
	}

	return nil
}

func (u *TableUsecase) DeleteTable(id uuid.UUID) error {
	table, err := u.tableRepository.GetTableByID(id)
	if err != nil {
		return fmt.Errorf("table not found: %v", err)
	}

	if err := u.tableRepository.DeleteTable(table); err != nil {
		return fmt.Errorf("failed to delete table: %v", err)
	}

	return nil
}

// Timeslot Usecase
func (u *TableUsecase) CreateTimeslot(request *dto.CreateTimeslotRequest) (*dto.TimeslotDetail, error) {
	// Validate and parse times
	startTime, err := u.parseTodayTime(request.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := u.parseTodayTime(request.EndTime)
	if err != nil {
		return nil, err
	}

	if !endTime.After(startTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Check for overlapping timeslots
	if overlap, err := u.isTimeslotOverlapping(startTime, endTime, uuid.Nil); err != nil {
		return nil, fmt.Errorf("error checking timeslot overlap: %v", err)
	} else if overlap {
		return nil, fmt.Errorf("timeslot overlaps with existing timeslots")
	}

	// Create timeslot
	timeslot := &models.Timeslot{
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := u.tableRepository.CreateTimeslot(timeslot); err != nil {
		return nil, err
	}

	// Create TableTimeslot for each table
	tables, err := u.tableRepository.GetAllTables()
	if err != nil {
		return nil, fmt.Errorf("failed to get tables for timeslot: %v", err)
	}

	for i, table := range tables {
		tableTimeslot := &models.TableTimeslot{
			TableID:       table.ID,
			TimeslotID:    timeslot.ID,
			Status:        "available", // ค่า default
			ReservedSeats: 0,
		}
		if err := u.tableRepository.CreateTableTimeslot(tableTimeslot); err != nil {
			return nil, fmt.Errorf("failed to create table timeslot at index %d: %v", i, err)
		}
	}


	return &dto.TimeslotDetail{
		ID:        timeslot.ID,
		StartTime: timeslot.StartTime.Format("15:04"),
		EndTime:   timeslot.EndTime.Format("15:04"),
	}, nil
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

func (u *TableUsecase) GetTimeslotByID(id uuid.UUID) (*dto.TimeslotDetail, error) {
	timeslot, err := u.tableRepository.GetTimeslotByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get timeslot: %v", err)
	}

	detail := &dto.TimeslotDetail{
		ID:        timeslot.ID,
		StartTime: timeslot.StartTime.Format("15:04"),
		EndTime:   timeslot.EndTime.Format("15:04"),
	}

	return detail, nil
}

func (u *TableUsecase) EditTimeslotDetails(id uuid.UUID, request *dto.EditTimeslotRequest) error {
	timeslot, err := u.tableRepository.GetTimeslotByID(id)
	if err != nil {
		return fmt.Errorf("timeslot not found: %v", err)
	}

	if request.StartTime != "" {
		timeslot.StartTime, err = u.parseTodayTime(request.StartTime)
		if err != nil {
			return err
		}
	}

	if request.EndTime != "" {
		timeslot.EndTime, err = u.parseTodayTime(request.EndTime)
		if err != nil {
			return err
		}
	}

	if !timeslot.EndTime.After(timeslot.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	// Check for overlapping timeslots excluding current timeslot
	if overlap, err := u.isTimeslotOverlapping(timeslot.StartTime, timeslot.EndTime, timeslot.ID); err != nil {
		return fmt.Errorf("error checking timeslot overlap: %v", err)
	} else if overlap {
		return fmt.Errorf("timeslot overlaps with existing timeslots")
	}

	if err := u.tableRepository.UpdateTimeslot(timeslot); err != nil {
		return fmt.Errorf("failed to update timeslot: %v", err)
	}

	return nil
}

func (u *TableUsecase) DeleteTimeslot(id uuid.UUID) error {

	timeslot, err := u.tableRepository.GetTimeslotByID(id)
	if err != nil {
		return fmt.Errorf("timeslot not found: %v", err)
	}

	// Delete all tableTimeslots under this timeslot
	if err := u.tableRepository.DeleteTableTimeslotByTimeslotID(timeslot.ID); err != nil {
		return fmt.Errorf("failed to delete related tableTimeslots: %v", err)
	}

	if err := u.tableRepository.DeleteTimeslot(timeslot); err != nil {
		return fmt.Errorf("failed to delete timeslot: %v", err)
	}

	return nil
}


// TableTimeslot Usecase
func (u *TableUsecase) GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) (*dto.TableTimeslotResponse, error) {

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
		// tableDetail := dto.TableDetail{
		// 	// ID:  table.ID,
		// 	TableRow: table.TableRow,
		// 	TableCol: table.TableCol,
		// 	MaxSeats: table.MaxSeats,
		// }

		detail := dto.TableTimeslotDetail{
			ID:             t.ID,
			TableRow: 		table.TableRow,
			TableCol: 		table.TableCol,
			MaxSeats: 		table.MaxSeats,
			Status:       	t.Status,
			ReservedSeats:  t.ReservedSeats,
			// MaxSeats:    	table.MaxSeats,
		}
		details = append(details, detail)
	}

	timeslot, err := u.tableRepository.GetTimeslotByID(timeslotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get timeslot: %v", err)
	}

	return &dto.TableTimeslotResponse{
		StartTime:   timeslot.StartTime.Format("15:04"),
		EndTime:     timeslot.EndTime.Format("15:04"),
		TableTimeslots: details,
	}, nil
}

func (u *TableUsecase) GetTableTimeslotByID(id uuid.UUID) (*dto.TableTimeslotDetail, error) {
	tableTimeslot, err := u.tableRepository.GetTableTimeslotByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get table timeslot: %v", err)
	}

	table, err := u.tableRepository.GetTableByID(tableTimeslot.TableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %v", err)
	}

	detail := &dto.TableTimeslotDetail{
		ID:             tableTimeslot.ID,
		Status:       	tableTimeslot.Status,
		ReservedSeats:  tableTimeslot.ReservedSeats,
		TableRow: 		table.TableRow,
		TableCol: 		table.TableCol,
		MaxSeats: 		table.MaxSeats,
	}

	return detail, nil
}

func (u *TableUsecase) GetNowTableTimeslots() (*dto.TableTimeslotResponse, error) {

	now := time.Now()

	timeslot, err := u.tableRepository.GetActiveTimeslot(now)
	if err != nil {
		return nil, fmt.Errorf("failed to get active timeslot: %v", err)
	}
	if timeslot == nil {
		return nil, fmt.Errorf("no active timeslot found")
	}

	response, err := u.GetTableTimeslotByTimeslotID(timeslot.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table timeslots for active timeslot: %v", err)
	}

	return response, nil
}