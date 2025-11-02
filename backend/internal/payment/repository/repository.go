package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/internal/db_model"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

/* -------------------- TRANSACTION -------------------- */

func (r *PaymentRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *PaymentRepository) GetAllTransactionsByUserID(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Where("user_id = ?", userID).Order("created_at ASC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

/* -------------------- PAYMENT METHOD -------------------- */

func (r *PaymentRepository) CreatePaymentMethod(method *models.PaymentMethod) error {
	return r.db.Create(method).Error
}

func (r *PaymentRepository) GetPaymentMethodByID(paymentMethodID uuid.UUID) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	if err := r.db.First(&method, "id = ?", paymentMethodID).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *PaymentRepository) GetPaymentMethodsByType(methodType string) ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	if methodType == "paid" {
		if err := r.db.Where("type IN ?", []string{methodType, "all"}).Find(&methods).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Where("type IN ?", []string{methodType, "all", "both"}).Find(&methods).Error; err != nil {
			return nil, err
		}
	}
	return methods, nil
}

func (r *PaymentRepository) GetAllPaymentMethods() ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	if err := r.db.Find(&methods).Error; err != nil {
		return nil, err
	}
	return methods, nil
}

/* -------------------- FOOD ORDER -------------------- */
func (r *PaymentRepository) GetFoodOrderByID(orderID uuid.UUID) (*models.FoodOrder, error) {
	var order models.FoodOrder
	err := r.db.First(&order, "id = ?", orderID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *PaymentRepository) GetFoodOrderByReservationID(reservationID uuid.UUID) (*models.FoodOrder, error) {
	var order models.FoodOrder
	err := r.db.Where("reservation_id = ?", reservationID).
		Order("created_at DESC").
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *PaymentRepository) UpdateFoodOrderStatus(orderID uuid.UUID, status string) error {
	return r.db.Model(&models.FoodOrder{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (r *PaymentRepository) GetTotalAmountForCustomerInOrder(orderID uuid.UUID, customerID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.Model(&models.FoodOrderItem{}).
		Where("food_order_id = ? AND customer_id = ?", orderID, customerID).
		Select("COALESCE(SUM(subtotal), 0)").
		Scan(&total).Error
	return total, err
}

func (r *PaymentRepository) GetTotalAmountByReservationID(reservationID uuid.UUID) (float64, error) {
    var total float64
    err := r.db.Model(&models.FoodOrder{}).
        Where("reservation_id = ?", reservationID).
        Select("COALESCE(SUM(total_amount), 0)").
        Scan(&total).Error
    return total, err
}


// ดึงร้านจาก FoodOrder
func (r *PaymentRepository) GetRestaurantByFoodOrderID(orderID uuid.UUID) (*models.Restaurant, error) {
    var restaurant models.Restaurant
    err := r.db.Table("restaurants").
        Joins("JOIN menu_items ON menu_items.restaurant_id = restaurants.id").
        Joins("JOIN food_order_items ON food_order_items.menu_item_id = menu_items.id").
        Where("food_order_items.food_order_id = ?", orderID).
        Limit(1).
        Scan(&restaurant).Error
    if err != nil {
        return nil, err
    }
    return &restaurant, nil
}


/* -------------------- TABLE RESERVATION -------------------- */

func (r *PaymentRepository) GetTableReservationMemberByCustomerID(reservationID uuid.UUID, customerID uuid.UUID) (*models.TableReservationMembers, error) {
	var member models.TableReservationMembers
	err := r.db.Where("reservation_id = ? AND customer_id = ?", reservationID, customerID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *PaymentRepository) GetAllMembersByTableReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error) {
	var members []models.TableReservationMembers
	if err := r.db.Where("reservation_id = ?", reservationID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *PaymentRepository) UpdateTableReservationMemberStatus(memberID uuid.UUID, status string) error {
	return r.db.Model(&models.TableReservationMembers{}).
		Where("id = ?", memberID).
		Update("status", status).Error
}

func (r *PaymentRepository) GetReservationByID(reservationID uuid.UUID) (*models.TableReservation, error) {
	var reservation models.TableReservation
	err := r.db.First(&reservation, "id = ?", reservationID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reservation, nil
}

func (r *PaymentRepository) UpdateTableReservationStatus(reservationID uuid.UUID, status string) error {
	return r.db.Model(&models.TableReservation{}).
		Where("id = ?", reservationID).
		Update("status", status).Error
}

/* -------------------- COMBINED UTILITY -------------------- */

// ใช้เมื่ออยากอัปเดตทุกอย่างใน transaction เดียว (atomic)
func (r *PaymentRepository) RunInTransaction(fn func(tx *gorm.DB) error) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	return tx.Commit().Error
}
