package interfaces

import (
	"context"

	"github.com/google/uuid"
	"backend/internal/notifications/dto"
	"gorm.io/gorm"
	db_model "backend/internal/db_model"

)

type NotiHandler interface {
	List(ctx context.Context, q dto.ListQuery) (*dto.ListResponse, error)
	MarkRead(ctx context.Context, id uuid.UUID, isRead bool) error
	MarkAllRead(ctx context.Context, receiverId uuid.UUID, receiverType string) (int, error)
	MockCreate(ctx context.Context, req dto.MockCreateRequest) (*dto.MockCreateResponse, error)
}

type ListFilter struct {
	ReceiverID   uuid.UUID
	ReceiverType string
	IsRead       *bool
	Offset       int
	Limit        int
	SortAsc      bool // true = created_at ASC, false = DESC
}

type NotiRepository interface {
    List(ctx context.Context, db *gorm.DB, f ListFilter) ([]db_model.Notifications, int64, error)
    CreateBulk(ctx context.Context, db *gorm.DB, notis []db_model.Notifications) error
    MarkRead(ctx context.Context, db *gorm.DB, id uuid.UUID, isRead bool) error
    MarkAllRead(ctx context.Context, db *gorm.DB, receiverId uuid.UUID, receiverType string) (int64, error)
}

type NotiUsecase interface {
	List(ctx context.Context, q dto.ListQuery) (*dto.ListResponse, error)
	MarkRead(ctx context.Context, id uuid.UUID, isRead bool) error
	MarkAllRead(ctx context.Context, receiverId uuid.UUID, receiverType string) (int, error)
	MockCreate(ctx context.Context, req dto.MockCreateRequest) (*dto.MockCreateResponse, error)
}