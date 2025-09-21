package usecase

import (
	"fmt"
	"strings"
	"mime/multipart"
	// "github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	// "context"

	"backend/internal/db_model"
	menuInterfaces "backend/internal/menu/interfaces"
	"backend/internal/restaurant/dto"
	"backend/internal/restaurant/interfaces"
	user "backend/internal/user/dto"
	"backend/internal/utils"
	"github.com/google/uuid"
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

	username := strings.TrimSpace(request.Username)

	// Create new restaurant
	restaurant := models.Restaurant{

		Username:     username,
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


func (u *RestaurantUsecase) EditInfo(restaurantID uuid.UUID, request *dto.EditRestaurantRequest) (*dto.EditRestaurantResponse, error) {
	// 1) ตรวจว่ามีร้านนี้ไหม
	if _, err := u.restaurantRepository.GetByID(restaurantID); err != nil {
		return nil, err
	}

	// 2) อัปเดต: repo ของคุณรวมทั้ง Name, MenuType และ AddOnMenuItem ไว้ใน "PartialUpdate"
	//    - ส่ง 4 อาร์กิวเมนต์ให้ครบ (uuid, *string, *string, []string)
	//    - ถ้าไม่อยากแก้ add-ons ให้ส่ง nil มาใน request.AddOnMenuItem ได้ (repo ควรตีความว่า "ไม่แตะ")
	if request.Name != nil || request.MenuType != nil || request.AddOnMenuItem != nil {
		_, err := u.restaurantRepository.PartialUpdate(restaurantID, request.Name, request.MenuType)
		if err != nil {
			return nil, err
		}
	}

	// 3) ดึงค่าล่าสุด
	updated, err := u.restaurantRepository.GetByID(restaurantID)
	if err != nil {
		return nil, err
	}

	// 4) map -> dto: จัดการ *string -> string ให้ตรงชนิดของ EditRestaurantResponse
	var nameVal string
	if updated.Name != nil {
		nameVal = *updated.Name
	}
	var menuTypeVal string
	if updated.MenuType != nil {
		menuTypeVal = *updated.MenuType
	}

	// 5) add-ons: เนื่องจากยังไม่มี GetAddOnItems ใน repo ของคุณ
	//    - ถ้ารอบนี้มีส่ง add-ons มาด้วย ให้ใช้ค่าที่เพิ่งแก้ (request.AddOnMenuItem)
	//    - ถ้าไม่ได้ส่งมา (nil) และอยากได้ค่าจริงจาก DB แนะนำเพิ่มเมธอด GetAddOnItems ภายหลัง
	addOns := []string{}
	if request.AddOnMenuItem != nil {
		addOns = request.AddOnMenuItem
	}

	resp := &dto.EditRestaurantResponse{
		ID:            updated.ID,     // ต้องเป็น uuid.UUID ใน DTO
		Name:          nameVal,        // from *string
		MenuType:      menuTypeVal,    // from *string
		AddOnMenuItem: addOns,         // จาก request รอบนี้ (ถ้า nil จะเป็น [] ว่าง)
	}
	return resp, nil
}
