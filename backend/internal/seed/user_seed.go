package seed

import (
	"fmt"
	"os"
	"path/filepath"
	"mime/multipart"
	"net/textproto"
	"gorm.io/gorm"
	"github.com/minio/minio-go/v7"

	"backend/internal/db_model"
	"backend/internal/utils"
)

func uploadSampleRestaurantImage(minioClient *minio.Client, filename string) (string, error) {
	basePath, _ := os.Getwd() // current working directory ของ process
	filePath := filepath.Join(basePath, "internal", "assets", "images", "restaurants", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file error: %v", err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileHeader := &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Size:     fileInfo.Size(),
		Header:   make(textproto.MIMEHeader),
	}
	fileHeader.Header.Set("Content-Type", "image/png")

	// upload to MinIO
	objectName := fmt.Sprintf("restaurants/%s", fileHeader.Filename)
	url, err := utils.UploadImage(file, fileHeader, "restaurant-pictures", objectName, minioClient)
	if err != nil {
		return "", err
	}
	return url, nil
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

func seedRestaurants(db *gorm.DB, minioClient *minio.Client) error {
	for i := 1; i <= 2; i++ {
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

		// ---- Upload image to MinIO ----
		filename := fmt.Sprintf("restaurant%02d.png", i)
		imgURL, err := uploadSampleRestaurantImage(minioClient, filename)
		if err != nil {
			return fmt.Errorf("failed to upload image for %s: %v", name, err)
		}

		// Create new restaurant
		restaurant := models.Restaurant{
			Username: username,
			Name:     name,
			Email:    email,
			Password: hashedPassword,
			WalletBalance: 200, // default wallet balance
			ProfilePic: &imgURL,
		}
		if err := db.Create(&restaurant).Error; err != nil {
			return err
		}
	}
	return nil
}