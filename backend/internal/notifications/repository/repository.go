package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	db_model "backend/internal/db_model"
	"backend/internal/notifications/interfaces"
)

type notiRepository struct {
	db *gorm.DB
}

func NewNotiRepository(db *gorm.DB) interfaces.NotiRepository {
	return &notiRepository{db: db}
}

func (r *notiRepository) List(ctx context.Context, db *gorm.DB, f interfaces.ListFilter) ([]db_model.Notifications, int64, error) {
	if db == nil {
		db = r.db
	}

	q := db.WithContext(ctx).Model(&db_model.Notifications{}).
		Where("receiver_id = ? AND receiver_type = ?", f.ReceiverID, f.ReceiverType)

	if f.IsRead != nil {
		q = q.Where("is_read = ?", *f.IsRead)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if f.SortAsc {
		q = q.Order("created_at ASC")
	} else {
		q = q.Order("created_at DESC")
	}

	var rows []db_model.Notifications
	if err := q.Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *notiRepository) CreateBulk(ctx context.Context, db *gorm.DB, notis []db_model.Notifications) error {
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Create(&notis).Error
}

func (r *notiRepository) MarkRead(ctx context.Context, db *gorm.DB, id uuid.UUID, isRead bool) error {
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).
		Model(&db_model.Notifications{}).
		Where("id = ?", id).
		Update("is_read", isRead).Error
}

func (r *notiRepository) MarkAllRead(ctx context.Context, db *gorm.DB, receiverId uuid.UUID, receiverType string) (int64, error) {
	if db == nil {
		db = r.db
	}
	tx := db.WithContext(ctx).
		Model(&db_model.Notifications{}).
		Where("receiver_id = ? AND receiver_type = ? AND is_read = ?", receiverId, receiverType, false).
		Update("is_read", true)

	return tx.RowsAffected, tx.Error
}
