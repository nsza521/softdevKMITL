package main

import (
	"fmt"
	// "log"
	// "os"
	"gorm.io/gorm"
	"backend/internal/db_model"
	"backend/internal/utils"
)

func InitAllSeedData(db *gorm.DB)  error {

	err := seedCustomers(db)
	if err != nil {
		return fmt.Errorf("error seeding customers: %v", err)
	}

	err = seedRestaurants(db)
	if err != nil {
		return fmt.Errorf("error seeding restaurants: %v", err)
	}

	// err = seedMenuItems(db)
	// if err != nil {
	// 	return fmt.Errorf("error seeding menu items: %v", err)
	// }

	// err = seedTables(db)
	// if err != nil {
	// 	return fmt.Errorf("error seeding tables: %v", err)
	// }

	// err = seedTimeSlots(db)
	// if err != nil {
	// 	return fmt.Errorf("error seeding time slots: %v", err)
	// }

	return nil
}

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

		// Check if restaurant already exists
		var count int64
		db.Model(&models.Restaurant{}).Where("username = ? OR email = ?", username, email).Count(&count)
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
			Email:    email,
			Password: hashedPassword,
		}
		if err := db.Create(&restaurant).Error; err != nil {
			return err
		}
	}
	return nil
}

// func seedMenuItems(db *gorm.DB) error {
// 	// Implement menu item seeding logic if needed
// 	return nil
// }

// func seedTables(db *gorm.DB) error {
// 	// Implement table seeding logic if needed
// 	return nil
// }

// func seedTimeSlots(db *gorm.DB) error {
// 	// Implement time slot seeding logic if needed
// 	return nil
// }