package interfaces

import (
	"backend/internal/db_model"
)

type CustomerRepository interface {
	GetByUsername(username string) (*models.Customer, error)
	Create(customer *models.Customer) error
	IsCustomerExists(username string, email string) (bool, error)
}