package usecase

import (
	"fmt"
	"time"
	"strings"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/minio/minio-go/v7"

	"backend/internal/customer/dto"
	"backend/internal/customer/interfaces"
	"backend/internal/utils"
	"backend/internal/db_model"
	// user "backend/internal/user/dto"
)

func stringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

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

func (u *CustomerUsecase) ListServedOrdersByCustomer(ctx context.Context, customerID string) ([]models.FoodOrder, error) {
	return u.customerRepository.ListServedOrdersByCustomer(ctx, customerID)
}

func (uc *CustomerUsecase) GetMyOrderHistory(c *gin.Context) ([]dto.OrderHistoryDay, error) {
	customerID := c.GetString("user_id")

	orders, err := uc.customerRepository.ListServedOrdersByCustomer(c, customerID)
    fmt.Println("GetMyHistory orders:", orders)

    if err != nil {
        return nil, err
    }

    // group orders by date (yyyy-mm-dd),
    // preserve order from newest to oldest
    dayBuckets := make(map[string]*dto.OrderHistoryDay)
    var dayOrderKeys []string

    for _, o := range orders {
        // สมมติ field ใน dbmodel.FoodOrder คือ OrderDate time.Time
        dayKey := o.OrderDate.Format("2006-01-02")

        if _, exists := dayBuckets[dayKey]; !exists {
            dayBuckets[dayKey] = &dto.OrderHistoryDay{
                Date:   dayKey,
                Orders: []dto.OrderHistoryItem{},
            }
            dayOrderKeys = append(dayOrderKeys, dayKey)
        }

        // map line items
        lineItems := make([]dto.OrderHistoryLineItem, 0, len(o.Items))
        for _, it := range o.Items {
            // map options
            opts := make([]dto.OrderHistoryLineItemOption, 0, len(it.Options))
            for _, op := range it.Options {
                opts = append(opts, dto.OrderHistoryLineItemOption{
                    OptionName: op.OptionName,
                    PriceDelta: op.PriceDelta,
                })
            }

            lineItems = append(lineItems, dto.OrderHistoryLineItem{
                MenuName:  it.MenuName,
                Quantity:  it.Quantity,
                UnitPrice: it.UnitPrice,
                Subtotal:  it.Subtotal,
                Options:   opts,
            })
        }
		// format order time 02-11-2025 18:11:00
        dayBuckets[dayKey].Orders = append(dayBuckets[dayKey].Orders, dto.OrderHistoryItem{
            OrderID:     o.ID.String(),
            Channel:     o.Channel,
            Note:        stringValue(o.Note),
            TotalAmount: o.TotalAmount,
            // OrderTime:   o.OrderDate.Format(time.RFC3339),
			OrderTime:   o.OrderDate.Format("02-01-2006 15:04:05"),
            Items:       lineItems,
        })
    }

    // flatten ตาม key ลำดับที่ append ไว้ (ใหม่ -> เก่า)
    result := make([]dto.OrderHistoryDay, 0, len(dayOrderKeys))
    for _, k := range dayOrderKeys {
        result = append(result, *dayBuckets[k])
    }

    fmt.Println("GetMyHistory result:", result)
    return result, nil
}