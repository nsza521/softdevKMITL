package usecase

import (
	"fmt"
	"time"
	"math"
	
	"github.com/google/uuid"
	"backend/internal/db_model"
	"backend/internal/reservation/dto"
	"backend/internal/reservation/interfaces"
	tableInterfaces "backend/internal/table/interfaces"
	customerInterfaces "backend/internal/customer/interfaces"
)

type TableReservationUsecase struct {
	tableReservationRepository interfaces.TableReservationRepository
	tableRepository             tableInterfaces.TableRepository
	customerRepository          customerInterfaces.CustomerRepository
}

func NewTableReservationUsecase(tableReservationRepository interfaces.TableReservationRepository, tableRepository tableInterfaces.TableRepository, customerRepository customerInterfaces.CustomerRepository) interfaces.TableReservationUsecase {
	return &TableReservationUsecase{
		tableReservationRepository: tableReservationRepository,
		tableRepository:             tableRepository,
		customerRepository:          customerRepository,
	}
}


// Table Reservation Usecase
func (u *TableReservationUsecase) getTableTimeslotStatus(reservedSeats int, maxSeats int) string {
	if reservedSeats >= int(0.8*float32(maxSeats)) {
		return "full"
	}
	if reservedSeats > 0 {
		return "partial"
	}
	return "available"
}

func (u *TableReservationUsecase) createNotRandomTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error) {
	// Implement logic for creating a non-random table reservation
	tableTimeslot, err := u.tableRepository.GetTableTimeslotByID(request.TableTimeslotID)
		if err != nil {
			return nil, err
		}

	status := tableTimeslot.Status
	if status == "full" || status == "expired" {
		return nil, fmt.Errorf("Table is not available for reservation")
	}

	// timeslot , err := u.tableRepository.GetTimeslotByID(tableTimeslot.TimeslotID)
	// if err != nil {
	// 	return err
	// }
	// if timeslot.EndTime.Before(request.ReserveTime) {
	// 	tableTimeslot.Status = "expired"
	// 	err = u.tableRepository.UpdateTableTimeslot(tableTimeslot)
	// 	if err != nil {	
	// 		return err
	// 	}
	// 	return fmt.Errorf("Timeslot is expired")
	// }

	table, err := u.tableRepository.GetTableByID(tableTimeslot.TableID)
	if err != nil {
		return nil, err
	}

	reservePeople := len(request.Members)
	if reservePeople <= 0 {
		return nil, fmt.Errorf("Reserve people must be greater than 0")
	}
	if reservePeople > table.MaxSeats {
		return nil, fmt.Errorf("Reserve people exceeds max seats of the table")
	}
	if reservePeople > (table.MaxSeats - tableTimeslot.ReservedSeats) {
		return nil, fmt.Errorf("Reserve people exceeds available seats of the table")
	}

	// Reserved people must not exceed 80% of max seats if not random
	random := request.Random
	if !random && reservePeople < int(0.8*float32(table.MaxSeats)) {
		random = true
	}

	reservation := &models.TableReservation{
		TableTimeslotID: request.TableTimeslotID,
		// CustomerID:      request.CustomerID,
		ReservePeople:   reservePeople,
		Random:          random,
		Status:          "pending", // Default status is "pending"
	}
	createdReservation, err := u.tableReservationRepository.CreateTableReservation(reservation)
	if err != nil {
		return nil, err
	}

	tableTimeslot.ReservedSeats += reservePeople
	if tableTimeslot.ReservedSeats > table.MaxSeats {
		return nil, fmt.Errorf("Reserved seats exceed max seats of the table")
	}

	tableTimeslot.Status = u.getTableTimeslotStatus(tableTimeslot.ReservedSeats, table.MaxSeats)

	err = u.tableRepository.UpdateTableTimeslot(tableTimeslot)
	if err != nil {	
		return nil, err
	}

	// Add members
	for _, member := range request.Members {
		err := u.CreateTableReservationMember(createdReservation.ID, member.Username)
		if err != nil {
			return nil, err
		}
	}

	timeslot , err := u.tableRepository.GetTimeslotByID(tableTimeslot.TimeslotID)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationDetail{
		CreateAt: 		 	 createdReservation.CreatedAt.Format("02-01-2006 15:04"),
		ReservationID:       createdReservation.ID,
		TableTimeslotID:     createdReservation.TableTimeslotID,
		ReservePeople:       createdReservation.ReservePeople,
		// Random:           	 createdReservation.Random,
		Status:          	 createdReservation.Status,
		Members:         	 request.Members,
		TableRow:	   		 table.TableRow,
		TableCol:	   		 table.TableCol,
		StartTime:   		 timeslot.StartTime.Format("15:04"),
		EndTime:     		 timeslot.EndTime.Format("15:04"),
	}, nil
}

func (u *TableReservationUsecase) createRandomTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error) {
	if request.Random == false {
		return nil, fmt.Errorf("Random must be true for random table reservation")
	}

	if len(request.Members) != 1 {
		return nil, fmt.Errorf("Random reservation allows only 1 member")
	}

	// 1. หา timeslot ที่ตรงกับเวลาปัจจุบัน
	now := time.Now()
	currentTimeslot, err := u.tableRepository.GetActiveTimeslot(now)
	if err != nil {
		return nil, err
	}
	if currentTimeslot == nil {
		return nil, fmt.Errorf("No available timeslot right now")
	}

	// 2. หา tabletimeslot ที่ยังว่าง

	availableTableTimeslot, err := u.tableRepository.GetAvailableTableTimeslot(currentTimeslot.ID)
	if err != nil {
		return nil, err
	}
	if availableTableTimeslot == nil {
		return nil, fmt.Errorf("No available tables in this timeslot")
	}

	// 4. สร้าง reservation
	reservation := &models.TableReservation{
		TableTimeslotID: availableTableTimeslot.ID,
		ReservePeople:   1,
		Random:          true,
		Status:          "pending",
	}
	createdReservation, err := u.tableReservationRepository.CreateTableReservation(reservation)
	if err != nil {
		return nil, err
	}

	// update reserved seats
	table, err := u.tableRepository.GetTableByID(availableTableTimeslot.TableID)
	if err != nil {
		return nil, err
	}
	minLeft := int(math.Ceil(0.8*float64(table.MaxSeats))) - availableTableTimeslot.ReservedSeats
	if minLeft < 0 {
		minLeft = 0
	}
	if minLeft+availableTableTimeslot.ReservedSeats > table.MaxSeats {
		return nil, fmt.Errorf("Reserved seats exceed max seats of the table")
	}
	availableTableTimeslot.ReservedSeats += 1
	availableTableTimeslot.Status = u.getTableTimeslotStatus(availableTableTimeslot.ReservedSeats, minLeft+availableTableTimeslot.ReservedSeats)
	if err := u.tableRepository.UpdateTableTimeslot(availableTableTimeslot); err != nil {
		return nil, err
	}

	// add member
	if err := u.CreateTableReservationMember(createdReservation.ID, request.Members[0].Username); err != nil {
		return nil, err
	}

	table, err = u.tableRepository.GetTableByID(availableTableTimeslot.TableID)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationDetail{
		CreateAt:       createdReservation.CreatedAt.Format("02-01-2006 15:04"),
		ReservationID:  createdReservation.ID,
		TableTimeslotID: createdReservation.TableTimeslotID,
		ReservePeople:  createdReservation.ReservePeople,
		Status:         createdReservation.Status,
		Members:        request.Members,
		TableRow:       table.TableRow,
		TableCol:       table.TableCol,
		StartTime:      currentTimeslot.StartTime.Format("15:04"),
		EndTime:        currentTimeslot.EndTime.Format("15:04"),
	}, nil
}


func (u *TableReservationUsecase) CreateTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error) {
	// เช็คว่า table ยังจองได้ไหม
	// เช็คว่า สมาชิกเกินจำนวนคนในโต๊ะไหม
	// ถ้าไม่เกินก็จองได้ ถ้าไม่ถึง 80% ของ maxSeats บังคับ random
	// ถ้าเกิน 80% ของ maxSeats ไม่ต้องบังคับ random
	// เช็คว่า customer มีสิทธิ์จองไหม
	// เช็คว่า customer มีการจองใน timeslot นี้อยู่แล้วไหม
	// เช็คว่า timeslot หมดอายุไหม
	// เช็คว่า tableTimeslot หมดอายุไหม
	// เช็คว่า tableTimeslot ยังว่างไหม
	// เช็คว่า reservePeople มากกว่า 0 ไหม
	// เช็คว่า reservePeople ไม่เกิน maxSeats ไหม
	// เช็คว่า reservePeople ไม่เกินจำนวนที่เหลือในโต๊ะไหม
	// ถ้า random = false ต้องเช็คว่า โต๊ะนี้ว่างไหม
	// ถ้า random = true ต้องเช็คว่า มีโต๊ะว่างไหม
	// ถ้าไม่มีโต๊ะว่าง ต้องเพิ่มเข้าคิว ??

	if request.TableTimeslotID != uuid.Nil {
		return u.createNotRandomTableReservation(request, customerID)
	} else {
		return u.createRandomTableReservation(request, customerID)
	}
}

func (u *TableReservationUsecase) GetTableReservationDetail(reservationID uuid.UUID, customerID uuid.UUID) (*dto.ReservationDetail, error) {
	err := u.isCustomerInReservation(reservationID, customerID)
	if err != nil {
		return nil, err
	}
	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return nil, err
	}

	members, err := u.GetAllMembersByReservationID(reservation.ID)
	if err != nil {
		return nil, err
	}
	membersDTO := []dto.Username{}
	for _, member := range members {
		customer, err := u.customerRepository.GetByID(member.CustomerID)
		if err != nil {
			return nil, err
		}
		membersDTO = append(membersDTO, dto.Username{Username: customer.Username})
	}

	tableTimeslot , err := u.tableRepository.GetTableTimeslotByID(reservation.TableTimeslotID)
	if err != nil {
		return nil, err
	}

	table , err := u.tableRepository.GetTableByID(tableTimeslot.TableID)
	if err != nil {
		return nil, err
	}

	timeslot , err := u.tableRepository.GetTimeslotByID(tableTimeslot.TimeslotID)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationDetail{
		CreateAt: 		 	 reservation.CreatedAt.Format("02-01-2006 15:04"),
		ReservationID:       reservation.ID,
		TableTimeslotID:     reservation.TableTimeslotID,
		ReservePeople:       reservation.ReservePeople,
		// Random:           	 reservation.Random,
		Status:          	 reservation.Status,
		Members:         	 membersDTO,
		TableRow:	   		 table.TableRow,
		TableCol:	   		 table.TableCol,
		StartTime:   		 timeslot.StartTime.Format("15:04"),
		EndTime:     		 timeslot.EndTime.Format("15:04"),
	}, nil
}

func (u *TableReservationUsecase) GetAllTableReservationHistory(customerID uuid.UUID) ([]dto.ReservationDetail, error) {
	reservationMembers, err := u.tableReservationRepository.GetAllReservationsByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	reservations := []dto.ReservationDetail{}
	for _, reservationMember := range reservationMembers {
		reservation, err := u.GetTableReservationDetail(reservationMember.ReservationID, customerID)
		if err != nil {
			return nil, err
		}

		if reservation.Status == "confirmed" {
			reservations = append(reservations, *reservation)
		}
	}
	return reservations, nil
}

func (u *TableReservationUsecase) DeleteTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error {
	err := u.isCustomerInReservation(reservationID, customerID)
	if err != nil {
		return err
	}

	err = u.tableReservationRepository.DeleteReservationMember(reservationID, customerID)
	if err != nil {
		return err
	}

	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil && err.Error() == "record not found" {
		return nil
	}
	if err != nil {
		return err
	}

	// Change TableReservation status to "cancelled" before deleting
	reservation.Status = "cancelled"
	err = u.tableReservationRepository.UpdateTableReservation(reservation)
	if err != nil {
		return err
	}
	err = u.tableReservationRepository.DeleteTableReservation(reservationID)
	if err != nil {
		return err
	}

	// Update TableTimeslot
	tableTimeslot, err := u.tableRepository.GetTableTimeslotByID(reservation.TableTimeslotID)
	if err != nil {
		return err
	}
	table, err := u.tableRepository.GetTableByID(tableTimeslot.TableID)
	if err != nil {
		return err
	}

	tableTimeslot.ReservedSeats -= reservation.ReservePeople
	if tableTimeslot.ReservedSeats < 0 {
		tableTimeslot.ReservedSeats = 0
	}

	tableTimeslot.Status = u.getTableTimeslotStatus(tableTimeslot.ReservedSeats, table.MaxSeats)

	return u.tableRepository.UpdateTableTimeslot(tableTimeslot)
}

// Table Reservation Members Usecase
func (u *TableReservationUsecase) isCustomerInReservation(reservationID uuid.UUID, customerID uuid.UUID) error {
	inReservation, err := u.tableReservationRepository.IsCustomerInReservation(reservationID, customerID)
	if err != nil {
		return err
	}
	if !inReservation {
		return fmt.Errorf("Customer is not in reservation")
	}
	return nil
}

func (u *TableReservationUsecase) isAllMembersConfirmed(reservationID uuid.UUID) bool {
	members, err := u.tableReservationRepository.GetAllMembersByReservationID(reservationID)
	if err != nil {
		return false
	}

	for _, member := range members {
		if member.Status != "confirmed" {
			return false
		}
	}
	return true
}

func (u *TableReservationUsecase) ConfirmTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error {
	err := u.isCustomerInReservation(reservationID, customerID)
	if err != nil {
		return err
	}

	reservationMember, err := u.tableReservationRepository.GetTableReservationMember(reservationID, customerID)
	if err != nil {
		return err
	}

	status := "confirmed"
	if reservationMember.Status == status {
		return fmt.Errorf("Customer has already confirmed the reservation")
	}

	reservationMember.Status = status
	err = u.tableReservationRepository.UpdateTableReservationMember(reservationMember)
	if err != nil {
		return err
	}

	if u.isAllMembersConfirmed(reservationID) {
		reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
		if err != nil {
			return err
		}

		reservation.Status = status
		err = u.tableReservationRepository.UpdateTableReservation(reservation)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *TableReservationUsecase) CreateTableReservationMember(reservationID uuid.UUID, username string) error {
	customer , err := u.customerRepository.GetByUsername(username)
	if err != nil {
		return err
	}
	member := &models.TableReservationMembers{
		ReservationID: reservationID,
		CustomerID:    customer.ID,
		Status:        "pending", // Default status is "pending"
	}
	return u.tableReservationRepository.CreateTableReservationMember(member)
}

func (u *TableReservationUsecase) CancelTableReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error {
	err := u.isCustomerInReservation(reservationID, customerID)
	if err != nil {
		return err
	}

	err = u.tableReservationRepository.DeleteReservationMember(reservationID, customerID)
	if err != nil {
		return err
	}

	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return err
	}
	if reservation.Status == "cancelled" {
		return nil
	}
	if reservation.Status == "confirmed" {
		return fmt.Errorf("Cannot cancel a confirmed reservation")
	}

	// Decrease reserve people by 1
	reservation.ReservePeople -= 1
	if reservation.ReservePeople < 0 {
		reservation.ReservePeople = 0
	}
	err = u.tableReservationRepository.UpdateTableReservation(reservation)
	if err != nil {
		return err
	}

	tableTimeslot, err := u.tableRepository.GetTableTimeslotByID(reservation.TableTimeslotID)
	if err != nil {
		return err
	}
	table, err := u.tableRepository.GetTableByID(tableTimeslot.TableID)
	if err != nil {
		return err
	}

	tableTimeslot.ReservedSeats -= 1
	if tableTimeslot.ReservedSeats < 0 {
		tableTimeslot.ReservedSeats = 0
	}

	tableTimeslot.Status = u.getTableTimeslotStatus(tableTimeslot.ReservedSeats, table.MaxSeats)

	err = u.tableRepository.UpdateTableTimeslot(tableTimeslot)
	if err != nil {	
		return err
	}

	if reservation.ReservePeople == 0 {
		err = u.tableReservationRepository.DeleteTableReservation(reservationID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *TableReservationUsecase) GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error) {
	return u.tableReservationRepository.GetAllMembersByReservationID(reservationID)
}

func (u *TableReservationUsecase) AddMemberToReservation(reservationID uuid.UUID, username string) (*dto.ReservationMemberDetail, error) {
	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return nil, err
	}

	members, err := u.GetAllMembersByReservationID(reservation.ID)
	if err != nil {
		return nil, err
	}

	if len(members) >= reservation.ReservePeople {
		return nil, fmt.Errorf("Reservation is full")
	}
	return nil, nil
}