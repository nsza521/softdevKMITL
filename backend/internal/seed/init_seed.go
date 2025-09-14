package seed

import (
	"fmt"

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

	if err := seedMenuTypesAndItems(db); err != nil {
		return fmt.Errorf("error seeding menu types & items: %v", err)
	}

	err = seedTableTimeslots(db)
	if err != nil {
		return fmt.Errorf("error seeding tables and time slots: %v", err)
	}

	return nil
}

