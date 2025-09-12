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