package dto

import (
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID			uuid.UUID	`json:"id"`
	Title		string		`json:"title"`
	Content		string		`json:"content"`
	Type		string		`json:"type"`	// e.g. "system", "order", "promo"
	ActionURL	*string		`json:"actionUrl"`
	ReceiverID   uuid.UUID `json:"receiverId"`
	ReceiverType string    `json:"receiverType"`
	IsRead       bool      `json:"isRead"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ListQuery struct {
	ReceiverID   uuid.UUID `form:"receiverId" binding:"required, stripuuid"`
	ReceiverType string    `form:"receiverType" binding:"required,oneof=customer restaurant"`
	IsRead       *bool     `form:"isRead"` // optional
	Page         int       `form:"page,default=1"`
	PageSize     int       `form:"pageSize,default=20"`
	Sort         string    `form:"sort,default=created_at_desc"` // created_at_desc|created_at_asc
}

type ListResponse struct {
	Items      []Notification `json:"items"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	TotalItems int64              `json:"totalItems"`
	TotalPages int                `json:"totalPages"`
}

type MarkReadRequest struct {
	IsRead bool `json:"isRead" binding:"required"`
}

type MockCreateRequest struct {
	ReceiverID   uuid.UUID `json:"receiverId" binding:"required"`
	ReceiverType string    `json:"receiverType" binding:"required,oneof=customer restaurant"`
	Count        int       `json:"count" binding:"omitempty,min=1,max=200"` // default 10
}

type MockCreateResponse struct {
	Inserted int `json:"inserted"`
}

type ListFilter struct {
    ReceiverID   uuid.UUID
    ReceiverType string
    IsRead       *bool
    Offset       int
    Limit        int
    SortAsc      bool
}