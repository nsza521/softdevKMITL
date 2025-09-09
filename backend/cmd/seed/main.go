package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	models "backend/internal/db_model"
	"backend/internal/seed"
)

func main() {
	_ = godotenv.Load()

	db, err := openMySQL()
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	// optional: ensure migration ก่อน
	if err := db.AutoMigrate(
		&models.Customer{},
		&models.Restaurant{},
		&models.Table{},
		&models.TableReservation{},
		&models.TableReservationMembers{},
		&models.Menu{},
		&models.MenuItem{},
		&models.FoodOrder{},
		&models.FoodOrderItem{},
		&models.Payment{},
		&models.Notifications{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	if err := seed.SeedDev(db); err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Println("✅ Seed completed")
}

func openMySQL() (*gorm.DB, error) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
