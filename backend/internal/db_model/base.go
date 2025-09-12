package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
    ID        uuid.UUID      `gorm:"type:char(36);primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
    if base.ID == "" {
        base.ID = uuid.New()// UUID v4
    }
    return
}
