package repository

import (
	"gorm.io/gorm"
)

type NotiRepository struct {
	db *gorm.DB
}

func NewNotiRepository(db *gorm.DB) *NotiRepository {
	return &NotiRepository{
		db: db,
	}
}