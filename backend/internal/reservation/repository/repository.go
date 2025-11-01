package repository

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"

	"backend/internal/db_model"
)

type TableReservationRepository struct {
	db *gorm.DB
}

func NewTableReservationRepository(db *gorm.DB) *TableReservationRepository {
	return &TableReservationRepository{
		db: db,
	}
}

// Table Reservation Repository
func (r *TableReservationRepository) CreateTableReservation(reservation *models.TableReservation) (*models.TableReservation, error) {
	if err := r.db.Create(reservation).Error; err != nil {
		return nil, err
	}
	
	createdReservation, err := r.GetTableReservationByID(reservation.ID)
	if err != nil {
		return nil, err
	}
	return createdReservation, nil
}

func (r *TableReservationRepository) GetTableReservationByID(id uuid.UUID) (*models.TableReservation, error) {
	var reservation models.TableReservation
	if err := r.db.First(&reservation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *TableReservationRepository) UpdateTableReservation(reservation *models.TableReservation) error {
	if err := r.db.Save(reservation).Error; err != nil {
		return err
	}
	return nil
}

func (r *TableReservationRepository) DeleteTableReservation(reservationID uuid.UUID) error {
	if err := r.db.Where("id = ?", reservationID).Delete(&models.TableReservation{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *TableReservationRepository) CountReservationsByCustomerAndDate(customerID uuid.UUID, date time.Time) (int64, error) {
    var count int64
	loc, _ := time.LoadLocation("Asia/Bangkok")
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
    endOfDay := startOfDay.Add(24 * time.Hour)

    err := r.db.
        Table("table_reservations").
        Joins("JOIN table_reservation_members m ON m.reservation_id = table_reservations.id").
        Where("m.customer_id = ? AND table_reservations.status IN (?) AND table_reservations.created_at BETWEEN ? AND ?",
            customerID,
            []string{"pending", "confirmed"},
            startOfDay,
            endOfDay).
        Count(&count).Error

    return count, err
}



// Table Reservation Members Repository
func (r *TableReservationRepository) CreateTableReservationMember(member *models.TableReservationMembers) error {
	return r.db.Create(member).Error
}

func (r *TableReservationRepository) GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error) {
	var members []models.TableReservationMembers
	if err := r.db.Where("reservation_id = ?", reservationID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *TableReservationRepository) IsCustomerInReservation(reservationID uuid.UUID, customerID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.TableReservationMembers{}).
		Where("reservation_id = ? AND customer_id = ?", reservationID, customerID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TableReservationRepository) GetAllReservationsByCustomerID(customerID uuid.UUID) ([]models.TableReservationMembers, error) {
	var reservations []models.TableReservationMembers
	if err := r.db.Where("customer_id = ?", customerID).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *TableReservationRepository) DeleteReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error {
	if err := r.db.Where("reservation_id = ? AND customer_id = ?", reservationID, customerID).Delete(&models.TableReservationMembers{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *TableReservationRepository) GetTableReservationMember(reservationID uuid.UUID, customerID uuid.UUID) (*models.TableReservationMembers, error) {
	var reservationMember models.TableReservationMembers
	if err := r.db.Where("reservation_id = ? AND customer_id = ?", reservationID, customerID).First(&reservationMember).Error; err != nil {
		return nil, err
	}
	return &reservationMember, nil
}

func (r *TableReservationRepository) UpdateTableReservationMember(member *models.TableReservationMembers) error {
	return r.db.Save(member).Error
}