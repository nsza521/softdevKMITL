package interfaces

import (
	"backend/internal/db_model"
)

type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) (*models.Restaurant, error)
	IsRestaurantExists(username string, email string) (bool, error)
	GetByUsername(username string) (*models.Restaurant, error)
	CreateBankAccount(bankAccount *models.BankAccount) error
}