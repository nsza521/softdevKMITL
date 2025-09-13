package seed

import (
	"fmt"
	"time"

	// "log"
	// "os"
	"backend/internal/db_model"
	"backend/internal/utils"

	"gorm.io/gorm"
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

	err = seedTableTimeslots(db)
	if err != nil {
		return fmt.Errorf("error seeding tables and time slots: %v", err)
	}

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

func seedTableTimeslots(db *gorm.DB) error {

	// Create tables
	var tables []models.Table
	for col := 1; col <= 3; col++ {
		for row := 1; row <= 6; row++ {
			table := models.Table{
				PeopleNum: 6,
				Row:      fmt.Sprintf("%c", 'A'+(row-1)),
				Col:      fmt.Sprintf("%d", col),
			}
			if err := db.Create(&table).Error; err != nil {
				return err
			}
			tables = append(tables, table)
		}
	}

	// Create time slots from 10:01 to 13:00 with 15-minute intervals
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	baseDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	gap := 1 * time.Minute
	start := baseDate.Add(10 * time.Hour + gap)   // 10:01 in Thailand
	end := baseDate.Add(13 * time.Hour)     // 13:00 in Thailand
	duration := 14 * time.Minute
	var timeSlots []models.Timeslot

	for t := start; t.Before(end); t = t.Add(duration + gap) {
		timeSlot := models.Timeslot{
			StartTime: t,
			EndTime:   t.Add(duration),
		}
		if err := db.Create(&timeSlot).Error; err != nil {
			return err
		}
		timeSlots = append(timeSlots, timeSlot)
	}


	// Create table-time slot associations
	for _, table := range tables {
		for _, timeSlot := range timeSlots {

			status := "available"
			if timeSlot.EndTime.In(loc).After(time.Now().In(loc)) {
				status = "expired"
			}

			tableTimeslot := models.TableTimeslot{
				TableID:    table.ID,
				TimeslotID: timeSlot.ID,
				Status:     status,
			}
			if err := db.Create(&tableTimeslot).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

