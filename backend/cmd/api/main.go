package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	models "backend/internal/db_model"
	"backend/internal/app"
)

func main() {
	// โหลด .env แค่ครั้งเดียว
	_ = godotenv.Load()

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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// AutoMigrate
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
		return nil, fmt.Errorf("error migrating database: %v", err)
	}
	log.Println("Database migrated successfully")
	return db, nil
}

func initMinIO() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
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

	return minioClient, nil
}
