package dto

import (
	"github.com/google/uuid"
)

type ProfileResponse struct {
	ID            	uuid.UUID   `json:"id"`
	Username  		string      `json:"username"`
	Email     		string      `json:"email"`
	FirstName 		string      `json:"first_name"`
	LastName  		string      `json:"last_name"`
	WalletBalance   float32     `json:"wallet_balance"`
}

type EditProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetFullnameRequest struct {
	Username string `json:"username" binding:"required"`
}

type GetFullnameResponse struct {
	Fullname string `json:"full_name"`
}