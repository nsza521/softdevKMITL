package models
import (
	"github.com/google/uuid"
)
type Notifications struct {
	Base
	Title			string		`gorm:"not null"`
	Content			string		`gorm:"not null"`
	ReceiverID 		uuid.UUID 	`gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReceiverType 	string 		`gorm:"not null"` // "customer" or "restaurant"
}