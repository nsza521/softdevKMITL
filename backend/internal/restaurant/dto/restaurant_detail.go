package dto

import (
	"github.com/google/uuid"
)

type RestaurantDetailResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	PictureURL *string   `json:"picture_url"`
	Email      string    `json:"email"`
	Status     string    `json:"status"`
}

type ChangeStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=open closed"`
}

// edit

type EditRestaurantRequest struct {
	Name          *string  `json:"name,omitempty"`
	MenuType      *string  `json:"menu_type,omitempty"`
	AddOnMenuItem []string `json:"add_on_menu_item"`
}

type EditRestaurantResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	MenuType      string    `json:"menu_type"`
	AddOnMenuItem []string  `json:"add_on_menu_item"`
}
