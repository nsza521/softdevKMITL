package models
import (
	"github.com/google/uuid"
	"time"
)

type NotificationType string

const (
	NotificationTypeSystem	NotificationType = "SYSTEM"
	NotificationTypeBooking	NotificationType = "BOOKING"
	NotificationTypePayment	NotificationType = "PAYMENT"
)

type Notifications struct {
	Base
	Title			string		`gorm:"not null"`
	Content			string		`gorm:"not null"`
	Type        	NotificationType		`gorm:"type:varchar(32);not null;default:'SYSTEM'"`
	ActionURL  		*string 	`gorm:"type:text"`
	ReceiverID 		uuid.UUID 	`gorm:"type:char(36);not null;index:idx_receiver"`
	ReceiverType 	string 		`gorm:"type:char(32);not null;index:idx_receiver"` // "customer" or "restaurant"
	IsRead			bool		`gorm:"not null;default:false;index"`
	CreatedAt		time.Time	`gorm:"index"`
}

func (Notifications) TableName() string { return "notifications" }