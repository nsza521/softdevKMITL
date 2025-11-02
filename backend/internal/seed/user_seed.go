package seed

import (
	"fmt"
	"gorm.io/gorm"

	"backend/internal/db_model"
	"backend/internal/utils"
)

func seedCustomers(db *gorm.DB) error {

	for i := 1; i <= 10; i++ {
		username := fmt.Sprintf("customer%02d", i)
		email := fmt.Sprintf("customer%02d@example.com", i)

		// Check if customer already exists
		var count int64
		db.Model(&models.Customer{}).Where("username = ? OR email = ?", username, email).Count(&count)
		if count > 0 {
			continue
		}
		
		// Hash password
		hashedPassword, err := utils.HashPassword("SecureP@ssw0rd")
		if err != nil {
			return err
		}

		// Create new customer
		customer := models.Customer{
			Username: username,
			Email:    email,
			Password: hashedPassword,
			FirstName: fmt.Sprintf("FirstName%02d", i),
			LastName:  fmt.Sprintf("LastName%02d", i),
		}
		if err := db.Create(&customer).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedRestaurants(db *gorm.DB) error {
	for i := 1; i <= 10; i++ {
		username := fmt.Sprintf("restaurant%02d", i)
		email := fmt.Sprintf("restaurant%02d@example.com", i)
		name := fmt.Sprintf("ร้านข้าว%02d", i)

		// Check if restaurant already exists
		var count int64
		db.Model(&models.Restaurant{}).Where("username = ? OR email = ? OR name = ?", username, email, name).Count(&count)
		if count > 0 {
			continue
		}
		
		// Hash password
		hashedPassword, err := utils.HashPassword("SecureP@ssw0rd")
		if err != nil {
			return err
		}

		// Create new restaurant
		restaurant := models.Restaurant{
			Username: username,
			Name:     name,
			Email:    email,
			Password: hashedPassword,
			WalletBalance: 200, // default wallet balance
		}
		if err := db.Create(&restaurant).Error; err != nil {
			return err
		}
	}
	return nil
}