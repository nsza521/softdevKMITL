package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"backend/internal/app"
	"backend/internal/db_model"
	"backend/internal/seed"
	"backend/internal/utils"
)

func main() {

	utils.BlacklistCleanup(5 * time.Minute)

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	db, err := initMySQL()
	if err != nil {
		log.Fatalf("Error initializing MySQL: %v", err)
	}

	minioClient, err := initMinIO()
	if err != nil {
		log.Fatalf("Error initializing MinIO: %v", err)
	}

	application := app.NewApp(db, minioClient)
	if err := application.Run(); err != nil {
		log.Fatalf("Error starting app: %v", err)
	}

}

func initMySQL() (*gorm.DB, error) {

	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FBangkok",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// AutoMigrate
	if err := db.AutoMigrate(
		&models.Customer{},
		&models.Restaurant{},
		&models.BankAccount{},
		&models.RestaurantAddOn{},
		&models.Table{},
		&models.Timeslot{},
		&models.TableTimeslot{},
		&models.TableReservation{},
		&models.TableReservationMembers{},
		&models.TopupHistory{},

		// Menu
		&models.MenuType{},
		&models.MenuItem{},
		&models.MenuTag{},
		&models.MenuAddOnGroup{},
		&models.MenuAddOnOption{},
		&models.MenuTypeAddOnGroup{},
		&models.MenuItemAddOnGroup{},

		&models.FoodOrder{},
		&models.FoodOrderItem{},
		&models.FoodOrderItemOption{},

		
		&models.PaymentMethod{},
		&models.Payment{},
		&models.Transaction{},
		&models.Notifications{},
	); err != nil {
		return nil, fmt.Errorf("error migrating database: %v", err)
	}
	log.Println("Database connected and migrated successfully")

	err = seed.InitAllSeedData(db)
	if err != nil {
		return nil, fmt.Errorf("error seeding database: %v", err)
	}

	log.Println("Database seeding completed successfully")

	return db, nil
}

func initMinIO() (*minio.Client, error) {

	endpoint := os.Getenv("MINIO_INTERNAL_ENDPOINT")
	if endpoint == "" {
		endpoint = "minio:9000" // default for docker compose
	}
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	// prevent trailing slash
	if len(endpoint) > 0 && endpoint[len(endpoint)-1] == '/' {
		endpoint = endpoint[:len(endpoint)-1]
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing MinIO client: %v", err)
	}
	log.Println("MinIO client initialized successfully")

	// Create buckets and sub-buckets
	roles := []string{"restaurant", "customer"}
	for _, role := range roles {
		bucketName, subBuckets := utils.GetBucketAndSubBuckets(role)
		err = utils.CreateBucketAndSubBuckets(minioClient, bucketName[0], subBuckets)
		if err != nil {
			return nil, fmt.Errorf("error setting up buckets for role %s: %v", role, err)
		}
		log.Printf("Buckets for role %s are set up", role)
	}

	return minioClient, nil
}
