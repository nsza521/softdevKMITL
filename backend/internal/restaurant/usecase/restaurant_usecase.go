package usecase

import (
	"fmt"
	"mime/multipart"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"backend/internal/restaurant/interfaces"
	"backend/internal/restaurant/dto"
	"backend/internal/utils"
	"backend/internal/db_model"
	user "backend/internal/user/dto"
	menuInterfaces "backend/internal/menu/interfaces"
)

type RestaurantUsecase struct {
	restaurantRepository interfaces.RestaurantRepository
	menuRepository       menuInterfaces.MenuRepository
	minioClient          *minio.Client
}

func NewRestaurantUsecase(restaurantRepository interfaces.RestaurantRepository, menuRepository menuInterfaces.MenuRepository, minioClient *minio.Client) interfaces.RestaurantUsecase {
	return &RestaurantUsecase{
		restaurantRepository: restaurantRepository,
		menuRepository:       menuRepository,
		minioClient:          minioClient,
	}
}

func (u *RestaurantUsecase) Register(request *dto.RegisterRestaurantRequest) error {

	// Check if restaurant exists
	exists, err := u.restaurantRepository.IsRestaurantExists(request.Username, request.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("restaurant already exists")
	}

	// Validate email format
	if !utils.IsValidEmail(request.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Check password strength
	if !utils.IsStrongPassword(request.Password) {
		return fmt.Errorf("password is not strong enough")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	// Create new restaurant
	restaurant := models.Restaurant{
		Username:     request.Username,
		Email:        request.Email,
		Password:     hashedPassword,
	}
	createdRestaurant, err := u.restaurantRepository.Create(&restaurant)
	if err != nil {
		return err
	}

	// Create bank account
	bankAccount := models.BankAccount{
		UserID:        createdRestaurant.ID,
		BankName:      request.BankName,
		AccountNumber: request.AccountNumber,
		AccountName:   request.AccountName,
	}
	err = u.restaurantRepository.CreateBankAccount(&bankAccount)
	if err != nil {
		return err
	}

	return nil
}

func (u *RestaurantUsecase) Login(request *user.LoginRequest) (string, error) {
	restaurant, err := u.restaurantRepository.GetByUsername(request.Username)
	if err != nil {
		return "", err
	}

	// Check password
	err = utils.VerifyPassword(request.Password, restaurant.Password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(restaurant.ID, restaurant.Username, "restaurant")
	if err != nil {
		return "", err
	}
	return token, nil

}

func (u *RestaurantUsecase) GetAll() ([]dto.RestaurantDetailResponse, error) {
	restaurants, err := u.restaurantRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var restaurantDetails []dto.RestaurantDetailResponse
	for _, r := range restaurants {
		detail := dto.RestaurantDetailResponse{
			ID:         r.ID,
			Username:   r.Username,
			PictureURL: r.ProfilePic,
			Email:      r.Email,
			Status:     r.Status,
		}
		restaurantDetails = append(restaurantDetails, detail)
	}

	return restaurantDetails, nil
}

func (u *RestaurantUsecase) UploadProfilePicture(restaurantID uuid.UUID, file *multipart.FileHeader) (string, error) {

	// Check if restaurant exists
	restaurant, err := u.restaurantRepository.GetByID(restaurantID)
	if err != nil {
		return "", err
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	// Upload to MinIO
	const bucketName = "restaurant-pictures"
	const subBucket = "restaurants"
	filename := restaurantID.String()
	objectName := fmt.Sprintf("%s/%s", subBucket, filename)

	url, err := utils.UploadImage(fileContent, file, bucketName, objectName, u.minioClient)
	if err != nil {
		return "", err
	}

	// Update restaurant profile picture URL
	if restaurant != nil {
		restaurant.ProfilePic = &url
	}
	err = u.restaurantRepository.Update(restaurant)
	if err != nil {
		return "", err
	}

	// presignURL, err := utils.GetPresignedURL(u.minioClient, bucketName, objectName)
	// if err != nil {
	// 	return "", err
	// }

	return url, nil
}

func (u *RestaurantUsecase) ChangeStatus(restaurantID uuid.UUID, request *dto.ChangeStatusRequest) error {
	// Check if restaurant exists
	restaurant, err := u.restaurantRepository.GetByID(restaurantID)
	if err != nil {
		return err
	}

	// Update status
	if request.Status != "" {
		restaurant.Status = request.Status
		return u.restaurantRepository.Update(restaurant)
	}

	return nil
}
