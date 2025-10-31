package dto

import (
	// "time"
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
	CreatedAt    string `json:"createdAt"`

	Attributes 	map[string]interface{} `json:"attributes,omitempty"`
}

type RequestReadAll struct {
	// ReceiverID uuid.UUID 	`json:"receiverId" binding:"required"`
	ReceiverType string    `json:"receiverType" binding:"required"`
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

type CreateEventRequest struct {
	Event        string      `json:"event" binding:"required"` // reserve_with | order_finished | order_canceled | reserve_success | reserve_failed
	ReceiverID   uuid.UUID   `json:"receiverId,omitempty"`
	ReceiverType string      `json:"receiverType" binding:"required,oneof=customer restaurant"`
	Data         interface{} `json:"data" binding:"required"` // payload เฉพาะแต่ละ event (struct ด้านล่าง)

	ReceiverUsername string `json:"receiverUsername,omitempty"`
}

type CreateEventResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// ---- payload เฉพาะเหตุการณ์ ----
type ReserveWithData struct {
	TableNo     string   `json:"tableNo"`
	When        string   `json:"when"`          // "19 ส.ค. 2025, เวลา xx:xx น."
	Restaurant  string   `json:"restaurant"`    // ชื่อร้าน
	Members     []string `json:"members"`       // ["Username", "Username", ...]
}

type OrderFinishedData struct {
	TableNo    string `json:"tableNo"`
	When       string `json:"when"`
	Restaurant string `json:"restaurant"`
	QueueNo    string `json:"queueNo,omitempty"`
}

type OrderCanceledData struct {
	TableNo    string `json:"tableNo"`
	When       string `json:"when"`
	Restaurant string `json:"restaurant"`
	Reason     string `json:"reason,omitempty"` // เช่น "ลูกค้ากดยกเลิก"
}

type ReserveSuccessData struct {
	TableNo    string `json:"tableNo"`
	When       string `json:"when"`
	Restaurant string `json:"restaurant"`
	Seat       int    `json:"seat"`
}

type ReserveFailedData struct {
	TableNo    string `json:"tableNo"`
	When       string `json:"when"`
	Restaurant string `json:"restaurant"`
}