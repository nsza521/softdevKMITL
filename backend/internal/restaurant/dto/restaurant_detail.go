package dto

import (
	"github.com/google/uuid"
)

type RestaurantDetailResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string  `json:"username"`
	PictureURL *string `json:"picture_url"`
	Email     string `json:"email"`
}

type EditRestaurantRequest struct {
	Username		*string `json:"username,omitempty"`
	Email			*string `json:"email,omitempty"`
	// PictureURL *string `json:"picture_url"`
	BankName    	*string `json:"bank_name,omitempty"`
	AccountNumber 	*string `json:"account_no,omitempty"`
	AccountName   	*string `json:"account_name,omitempty"`
	// PictureURL   *string `json:"picture_url"`
}