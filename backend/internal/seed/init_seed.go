package seed

import (
	"fmt"

	"gorm.io/gorm"
)

func InitAllSeedData(db *gorm.DB)  error {

	if err := seedCustomers(db); err != nil {
		return fmt.Errorf("error seeding customers: %v", err)
	}

	if err := seedRestaurants(db); err != nil {
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

	if err := seedFixedForNoodleShop(db); err != nil {
		return fmt.Errorf("error seeding noodle shop: %v", err)
	}

	return nil
}

