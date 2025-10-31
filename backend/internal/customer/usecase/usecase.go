package usecase

import (
	"fmt"
	"time"
	"strings"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/minio/minio-go/v7"

	"backend/internal/customer/dto"
	"backend/internal/customer/interfaces"
	"backend/internal/utils"
	"backend/internal/db_model"
	// user "backend/internal/user/dto"
)

type CustomerUsecase struct {
	customerRepository interfaces.CustomerRepository
	minioClient       *minio.Client
}

func NewCustomerUsecase(customerRepository interfaces.CustomerRepository, minioClient *minio.Client) interfaces.CustomerUsecase {
	return &CustomerUsecase{
		customerRepository: customerRepository,
		minioClient:       minioClient,
	}
}

func (u *CustomerUsecase) Register(request *dto.RegisterCustomerRequest) error {
	
	// Check if customer exists
	exists, err := u.customerRepository.IsCustomerExists(request.Username, request.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("customer already exists")
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

	firstName := strings.TrimSpace(strings.ToLower(request.FirstName))
	lastName := strings.TrimSpace(strings.ToLower(request.LastName))
	username := strings.TrimSpace(request.Username)

	// Create new customer
	customer := models.Customer{
		Username:     username,
		Email:        request.Email,
		FirstName:    firstName,
		LastName:     lastName,
		Password:     hashedPassword,
	}

	return  u.customerRepository.Create(&customer)
}

func (u *CustomerUsecase) Login(request *dto.LoginRequest) (string, error) {
	customer, err := u.customerRepository.GetByUsername(request.Username)
	if err != nil {
		return "", err
	}

	// hashedPassword, err := utils.HashPassword(request.Password)
	// if err != nil {
	// 	return "", err
	// }

	// Check password
	err = utils.VerifyPassword(request.Password, customer.Password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(customer.ID, customer.Username, "customer")
	if err != nil {
		return "", err
	}
	return token, nil

}

func (u *CustomerUsecase) Logout(token string, expiry time.Time) error {
	utils.BlacklistToken(token, expiry.Unix())
	return nil
}

func (u *CustomerUsecase) GetProfile(customerID uuid.UUID) (*dto.ProfileResponse, error) {
	customer, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return nil, err
	}

	response := &dto.ProfileResponse{
		ID:        customer.ID,
		Username:  customer.Username,
		Email:     customer.Email,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		WalletBalance: customer.WalletBalance,
	}
	return response, nil
}

func (u *CustomerUsecase) EditProfile(customerID uuid.UUID, request *dto.EditProfileRequest) error {
	customer, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return err
	}

	// Update fields
	if request.FirstName != "" {
		customer.FirstName = request.FirstName
	}
	if request.LastName != "" {
		customer.LastName = request.LastName
	}
	// if request.Email != "" {
	// 	// Validate email format
	// 	if !utils.IsValidEmail(request.Email) {
	// 		return fmt.Errorf("invalid email format")
	// 	}
	// 	customer.Email = request.Email
	// }

	return u.customerRepository.Update(customer)
}

func (u *CustomerUsecase) GetFullnameByUsername(customerID uuid.UUID, request *dto.GetFullnameRequest) (*dto.GetFullnameResponse, error) {
	customer, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	customer, err = u.customerRepository.GetByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	name, err := utils.ToTitleCase(customer.FirstName, customer.LastName)
	if err != nil {
		return nil, err
	}
	fullName := &dto.GetFullnameResponse{
		Fullname: name,
	}

	return fullName, nil
}

func (u *CustomerUsecase) GetFirstnameByUsername(customerID uuid.UUID, request *dto.GetFullnameRequest) (*dto.GetFirstnameResponse, error) {
	customer, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	customer, err = u.customerRepository.GetByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	firstName := &dto.GetFirstnameResponse{
		FirstName: customer.FirstName,
	}

	return firstName, nil
}

func (u *CustomerUsecase) GetQRCode(customerID uuid.UUID) (string, error) {
	customer, err := u.customerRepository.GetByID(customerID)
	if err != nil {
		return "", err
	}
	if customer == nil {
		return "", fmt.Errorf("customer not found")
	}

	// prepare data for QR code
	qrData := fmt.Sprintf("username: %s", customer.Username)
	
	// generate QR code
	qr, err := qrcode.New(qrData, qrcode.Medium)
	if err != nil {
		return "", err
	}
	qr.DisableBorder = false
	
	const size = 256
	png, err := qr.PNG(size)
	if err != nil {
		return "", err
	}
	url, err := utils.UploadBytes(
		png, "customer-pictures", 
		fmt.Sprintf("qr-codes/%s.png", customer.ID), 
		u.minioClient, 
		"image/png")

	return url, nil
}
