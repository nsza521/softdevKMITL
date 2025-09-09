package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
    // ID        uuid.UUID      `gorm:"type:binary(16);primaryKey"`
    ID        string         `gorm:"primaryKey;type:char(36)"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
    if base.ID == "" {
        base.ID = uuid.New().String() // UUID v4
    }
    return
}
