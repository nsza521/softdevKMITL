package interfaces

import (
	"github.com/google/uuid"
	"backend/internal/db_model"
)

type CustomerRepository interface {
	GetByUsername(username string) (*models.Customer, error)
	Create(customer *models.Customer) error
	IsCustomerExists(username string, email string) (bool, error)
	GetByID(id uuid.UUID) (*models.Customer, error)
	Update(customer *models.Customer) error
}