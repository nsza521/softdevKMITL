package usecase

import (
	"fmt"
	"time"
	
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
func (u *TableReservationUsecase) getTableTimeslotStatus(reservedSeats int, maxSeats int, random bool) string {
	if !random && reservedSeats >= int(0.8*float32(maxSeats)) {
		return "full"
	}
	if random && reservedSeats >= maxSeats {
		return "full"
	}
	if reservedSeats > 0 {
		return "partial"
	}
	return "available"
}

func (u *TableReservationUsecase) CreateNotRandomTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error) {
	// today := time.Now()
	// count, err := u.tableReservationRepository.CountReservationsByCustomerAndDate(customerID, today)
	// if err != nil {
	// 	return nil, err
	// }
	// if count >= 2 {
	// 	return nil, fmt.Errorf("You have reached the daily reservation limit (2 per day)")
	// }

	tableTimeslot, err := u.tableRepository.GetTableTimeslotByID(request.TableTimeslotID)
		if err != nil {
			return nil, err
		}

	status := tableTimeslot.Status
	if status == "full" || status == "expired" {
		return nil, fmt.Errorf("Table is not available for reservation")
	}
	if status == "partial" {
		return nil,  fmt.Errorf("TableTimeslot is already reserved")
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

	tableTimeslot.Status = u.getTableTimeslotStatus(tableTimeslot.ReservedSeats, table.MaxSeats, random)
	err = u.tableRepository.UpdateTableTimeslot(tableTimeslot)
	if err != nil {	
		return nil, err
	}

	// Add members
	for _, member := range request.Members {
		err := u.CreateTableReservationMember(createdReservation.ID, member.Username, "")
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

func (u *TableReservationUsecase) CreateRandomTableReservation(request *dto.CreateRandomTableReservationRequest, customerID uuid.UUID) (*dto.RandomReservationDetail, error) {
	customer, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return nil, err
	}

	// find current active timeslot
	// now := time.Now()
	// currentTimeslot, err := u.tableRepository.GetActiveTimeslot(now)
	// if err != nil {
	// 	return nil, err
	// }
	// if currentTimeslot == nil {
	// 	return nil, fmt.Errorf("No available timeslot right now")
	// }

	// timeslot, err := u.tableRepository.GetTimeslotByID(request.TimeslotID)
	// if err != nil {
	// 	return nil, err
	// }

	// today := time.Now()
	// count, err := u.tableReservationRepository.CountReservationsByCustomerAndDate(customerID, today)
	// if err != nil {
	// 	return nil, err
	// }
	// if count >= 2 {
	// 	return nil, fmt.Errorf("You have reached the daily reservation limit (2 per day)")
	// }

	// find available tableTimeslot
	availableTableTimeslots, err := u.tableRepository.GetAllAvailableTableTimeslot(request.TimeslotID)
	if err != nil {
		return nil, err
	}
	if len(availableTableTimeslots) == 0 {
		return nil, fmt.Errorf("No available tables in this timeslot")
	}

	for _, tableTimeslot := range availableTableTimeslots {
		reservations, err := u.tableReservationRepository.GetAllTableReservationByTableTimeslotID(tableTimeslot.ID)
		if err != nil {
			return nil, err
		}

		alreadyReserved := false
		for _, reservation := range reservations {
			err = u.isCustomerInReservation(reservation.ID, customerID)
			if err == nil {
				alreadyReserved = true
				break
			}
		}
		if alreadyReserved {
			continue
		}

		reservation := &models.TableReservation{
			TableTimeslotID: tableTimeslot.ID,
			ReservePeople:   1,
			Random:          true,
			Status:          "pending",
		}

		createdReservation, err := u.tableReservationRepository.CreateTableReservation(reservation)
		if err != nil {
			return nil, err
		}

		table, err := u.tableRepository.GetTableByID(tableTimeslot.TableID)
		if err != nil {
			return nil, err
		}

		tableTimeslot.ReservedSeats += 1
		tableTimeslot.Status = u.getTableTimeslotStatus(tableTimeslot.ReservedSeats, table.MaxSeats, true)
		if err := u.tableRepository.UpdateTableTimeslot(&tableTimeslot); err != nil {
			return nil, err
		}

		if err := u.CreateTableReservationMember(createdReservation.ID, customer.Username, "confirmed"); err != nil {
			return nil, err
		}

		return &dto.RandomReservationDetail{
			ReservationID:   createdReservation.ID,
			TableTimeslotID: createdReservation.TableTimeslotID,
		}, nil
	}
	// member := dto.Username{Username: customer.Username}
	return nil, fmt.Errorf("No available tables for reservation")
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
		return u.CreateNotRandomTableReservation(request, customerID)
	} else {
		return nil, fmt.Errorf("Not implemented yet")
		// return u.CreateRandomTableReservation(, customerID)
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

func (u *TableReservationUsecase) GetTableReservationOwnerDetail(reservationID uuid.UUID) (*dto.OwnerDetail, error) {
	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return nil, err
	}

	reservations, err := u.tableReservationRepository.GetAllTableReservationByTableTimeslotID(reservation.TableTimeslotID)
	if err != nil {
		return nil, err
	}
	var tableMembers []models.TableReservationMembers
	for _, res := range reservations {
		members, err := u.tableReservationRepository.GetAllMembersByReservationID(res.ID)
		if err != nil {
			return nil, err
		}
		tableMembers = append(tableMembers, members...)
	}

	ownerMember := tableMembers[0]
	customer, err := u.customerRepository.GetByID(ownerMember.CustomerID)
	if err != nil {
		return nil, err
	}

	return &dto.OwnerDetail{
		TableTimeslotID: reservation.TableTimeslotID,
		OwnerUsername:   customer.Username,
		OwnerFirstname:  customer.FirstName,
	}, nil
}

func (u *TableReservationUsecase) GetAllTableReservationHistory(customerID uuid.UUID) ([]dto.ReservationDetail, error) {
	reservationMembers, err := u.tableReservationRepository.GetAllTableReservationsByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	reservations := []dto.ReservationDetail{}
	status := "completed"
	for _, reservationMember := range reservationMembers {
		if reservationMember.Status != status {
			continue
		}
		
		reservation, err := u.GetTableReservationDetail(reservationMember.ReservationID, customerID)
		if err != nil {
			return nil, err
		}

		if reservation.Status == status {
			reservations = append(reservations, *reservation)
		}
	}
	return reservations, nil
}

func (u *TableReservationUsecase) GetAllTableReservationByCustomerID(customerID uuid.UUID) ([]dto.ReservationDetail, error) {
	reservationMembers, err := u.tableReservationRepository.GetAllTableReservationsByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	reservations := []dto.ReservationDetail{}
	for _, reservationMember := range reservationMembers {
		reservation, err := u.GetTableReservationDetail(reservationMember.ReservationID, customerID)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, *reservation)
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

	tableTimeslot.Status = u.getTableTimeslotStatus(tableTimeslot.ReservedSeats, table.MaxSeats, false)

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

func (u *TableReservationUsecase) getAllMemberByTableTimeslotID(tableTimeslotID uuid.UUID) ([]models.TableReservationMembers, error) {
	reservations, err := u.tableReservationRepository.GetAllTableReservationByTableTimeslotID(tableTimeslotID)
	if err != nil {
		return nil, err
	}	

	var allMembers []models.TableReservationMembers
	for _, reservation := range reservations {
		members, err := u.tableReservationRepository.GetAllMembersByReservationID(reservation.ID)
		if err != nil {
			return nil, err
		}
		allMembers = append(allMembers, members...)
	}
	return allMembers, nil
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

func (u *TableReservationUsecase) ConfirmMemberInTableReservation(reservationID uuid.UUID, customerID uuid.UUID) (*dto.ConfirmedStatusDetail, error) {
	err := u.isCustomerInReservation(reservationID, customerID)
	if err != nil {
		return nil, err
	}

	reservationMember, err := u.tableReservationRepository.GetTableReservationMember(reservationID, customerID)
	if err != nil {
		return nil, err
	}

	status := "confirmed"
	if reservationMember.Status == status {
		return nil, fmt.Errorf("Customer has already confirmed the reservation")
	}

	reservationMember.Status = status
	err = u.tableReservationRepository.UpdateTableReservationMember(reservationMember)
	if err != nil {
		return nil, err
	}

	members, err := u.tableReservationRepository.GetAllMembersByReservationID(reservationID)
	if err != nil {
		return nil, err
	}

	memberDetails := []dto.MemberStatus{}
	for _, member := range members {
		customer, err := u.customerRepository.GetByID(member.CustomerID)
		if err != nil {
			return nil, err
		}
		memberDetails = append(memberDetails, dto.MemberStatus{
			Username: customer.Username,
			Status:   member.Status,
		})
	}

	return &dto.ConfirmedStatusDetail{
		ReservationID:  reservationID,
		Members:       memberDetails,
	}, nil
}

func (u *TableReservationUsecase) ConfirmTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error {
	err := u.isCustomerInReservation(reservationID, customerID)
	if err != nil {
		return err
	}

	// reservationMember, err := u.tableReservationRepository.GetTableReservationMember(reservationID, customerID)
	// if err != nil {
	// 	return err
	// }

	status := "completed"
	members, err := u.tableReservationRepository.GetAllMembersByReservationID(reservationID)
	if err != nil {
		return err
	}

	for _, member := range members {
		member.Status = status
		err = u.tableReservationRepository.UpdateTableReservationMember(&member)
		if err != nil {
			return err
		}
	}

	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return err
	}

	reservation.Status = status
	err = u.tableReservationRepository.UpdateTableReservation(reservation)
	if err != nil {
		return err
	}

	// if reservationMember.Status == status {
	// 	return fmt.Errorf("Customer has already confirmed the reservation")
	// }

	// reservationMember.Status = status
	// err = u.tableReservationRepository.UpdateTableReservationMember(reservationMember)
	// if err != nil {
	// 	return err
	// }

	// if u.isAllMembersConfirmed(reservationID) {
	// 	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	reservation.Status = status
	// 	err = u.tableReservationRepository.UpdateTableReservation(reservation)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (u *TableReservationUsecase) GetTableReservationStatus(reservationID uuid.UUID, customerID uuid.UUID) (*dto.ReservationStatusDetail, error) {
	err := u.isCustomerInReservation(reservationID, customerID)
	if err != nil {
		return nil, err
	}

	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return nil, err
	}

	members, err := u.tableReservationRepository.GetAllMembersByReservationID(reservationID)
	if err != nil {
		return nil, err
	}

	var paidMembersCount int = 0
	for _, member := range members {
		if member.Status == "paid" || member.Status == "paid_pending" {
			paidMembersCount++
		}
	}

	return &dto.ReservationStatusDetail{
		ReservationStatus:   reservation.Status,
		TotalPeople:         reservation.ReservePeople,
		ConfirmedPaidPeople: paidMembersCount,
	}, nil
}

func (u *TableReservationUsecase) GetTableReservationTimeRemaining(reservationID uuid.UUID, customerID uuid.UUID) (*dto.ReservationTime, error) {

	if err := u.isCustomerInReservation(reservationID, customerID); err != nil {
		return nil, err
	}

	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return nil, err
	}

	// now := time.Now()

	// Expiration time is 5 minutes after creation
	expirationTime := reservation.CreatedAt.Add(5 * time.Minute)
	timeRemaining := time.Until(expirationTime)
	timeout := timeRemaining <= 0

	var formattedRemaining string
	if timeout {
		formattedRemaining = "00:00"
	} else {
		minutes := int(timeRemaining.Minutes())
		seconds := int(timeRemaining.Seconds()) % 60
		formattedRemaining = fmt.Sprintf("%02d:%02d", minutes, seconds)
	}

	return &dto.ReservationTime{
		TimeRemaining: formattedRemaining,
		Timeout:       timeout,
	}, nil
}


func (u *TableReservationUsecase) CreateTableReservationMember(reservationID uuid.UUID, username string, status string) error {
	customer , err := u.customerRepository.GetByUsername(username)
	if err != nil {
		return err
	}

	if status == "" {
		status = "pending"
	}
	member := &models.TableReservationMembers{
		ReservationID: reservationID,
		CustomerID:    customer.ID,
		Status:        status,
	}
	return u.tableReservationRepository.CreateTableReservationMember(member)
}

func (u *TableReservationUsecase) CancelTableReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error {
	if err := u.isCustomerInReservation(reservationID, customerID); err != nil {
		return err
	}

	if err := u.tableReservationRepository.DeleteReservationMember(reservationID, customerID); err != nil {
		return err
	}

	reservation, err := u.tableReservationRepository.GetTableReservationByID(reservationID)
	if err != nil {
		return err
	}

	// Check if reservation is already cancelled
	if reservation.Status == "cancelled" {
		return nil
	}
	if reservation.Status == "confirmed" {
		return fmt.Errorf("Cannot cancel a confirmed reservation")
	}

	reservation.ReservePeople -= 1
	if reservation.ReservePeople < 0 {
		reservation.ReservePeople = 0
	}

	// Update reservation information
	if err := u.tableReservationRepository.UpdateTableReservation(reservation); err != nil {
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

	tableTimeslot.ReservedSeats -= 1
	if tableTimeslot.ReservedSeats < 0 {
		tableTimeslot.ReservedSeats = 0
	}

	tableTimeslot.Status = u.getTableTimeslotStatus(tableTimeslot.ReservedSeats, table.MaxSeats, false)
	if err := u.tableRepository.UpdateTableTimeslot(tableTimeslot); err != nil {
		return err
	}

	// Get all members of the reservation
	members, err := u.tableReservationRepository.GetAllMembersByReservationID(reservationID)
	if err != nil {
		return err
	}

	// If there are no remaining members → delete the entire reservation
	if len(members) == 0 {
		if err := u.tableReservationRepository.DeleteTableReservation(reservationID); err != nil {
			return err
		}
		return nil
	}

	// If all members have paid → update reservation status to "paid"
	allPaid := true
	for _, member := range members {
		if member.Status != "paid" {
			allPaid = false
			break
		}
	}
	if allPaid {
		reservation.Status = "paid"
		if err := u.tableReservationRepository.UpdateTableReservation(reservation); err != nil {
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