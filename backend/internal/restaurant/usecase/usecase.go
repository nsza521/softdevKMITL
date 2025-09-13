package usecase

import (
	"fmt"

	"context"

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
}

func NewRestaurantUsecase(restaurantRepository interfaces.RestaurantRepository, menuRepository menuInterfaces.MenuRepository) interfaces.RestaurantUsecase {
	return &RestaurantUsecase{
		restaurantRepository: restaurantRepository,
		menuRepository:       menuRepository,
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
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
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
		}
		restaurantDetails = append(restaurantDetails, detail)
	}

	return restaurantDetails, nil
}

func (u *RestaurantUsecase) EditRestaurant(
	ctx context.Context,
	id , // path param จาก handler
	userID , role string, // ออกมาจาก JWT/middleware
	req dto.EditRestaurantRequest, // body JSON (pointer fields)
) (dto.RestaurantDetailResponse, error) {

	rid, err := uuid.Parse(id)
	if err != nil {
		return dto.RestaurantDetailResponse{}, fmt.Errorf("invalid restaurant id")
	}
	current, err := u.restaurantRepository.GetByID(rid)
	if err != nil {
		return dto.RestaurantDetailResponse{}, err
	}

	if current.ID.String() != userID {
		return dto.RestaurantDetailResponse{}, fmt.Errorf("forbidden")
	}

	changes := map[string]any{}
	if req.Username != nil {
		changes["username"] = *req.Username
	}
	if req.Email != nil {
		changes["email"] = *req.Email
	}
	if req.BankName != nil {
		changes["bank_name"] = *req.BankName
	}
	if req.AccountNumber != nil {
		changes["account_number"] = *req.AccountNumber
	}
	if req.AccountName != nil {
		changes["account_name"] = *req.AccountName
	}

	if req.Email != nil && !utils.IsValidEmail(*req.Email) {
		return dto.RestaurantDetailResponse{}, fmt.Errorf("invalid email format")
	}

	if len(changes) == 0 {
		return dto.RestaurantDetailResponse{
			ID:         current.ID,
			Username:   current.Username,
			PictureURL: current.ProfilePic,
			Email:      current.Email,
		}, nil
	}

	// changes["updated_at"] = time.Now()

	if err := u.restaurantRepository.PartialUpdate(ctx, id, changes); err != nil {
		return dto.RestaurantDetailResponse{}, err
	}
	updated, err := u.restaurantRepository.GetByID(rid)
	if err != nil { return dto.RestaurantDetailResponse{}, err }

	resp := dto.RestaurantDetailResponse{
		ID:         updated.ID,
		Username:   updated.Username,
		PictureURL: updated.ProfilePic, // ให้ตรงกับ db_model.Restaurant
		Email:      updated.Email,
	}
	return resp, nil
}
