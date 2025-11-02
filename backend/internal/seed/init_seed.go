package seed

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/minio/minio-go/v7"
)

func InitAllSeedData(db *gorm.DB, minioClient *minio.Client) error {

	if err := seedCustomers(db); err != nil {
		return fmt.Errorf("error seeding customers: %v", err)
	}

	if err := seedRestaurants(db, minioClient); err != nil {
		return fmt.Errorf("error seeding restaurants: %v", err)
	}

	if err := seedMenuTypesAndItems(db); err != nil {
		return fmt.Errorf("error seeding menu types & items: %v", err)
	}

	if err := seedTableTimeslots(db); err != nil {
		return fmt.Errorf("error seeding tables and time slots: %v", err)
	}

	if err := seedPaymentMethods(db); err != nil {
		return fmt.Errorf("error seeding payment methods: %v", err)
	}

	if err := seedFixedForNoodleShop(db, minioClient); err != nil {
		return fmt.Errorf("error seeding noodle shop: %v", err)
	}

	if err := RunSeedOrders(db); err != nil {
		return fmt.Errorf("error seeding orders: %v", err)
	}

	return nil
}

