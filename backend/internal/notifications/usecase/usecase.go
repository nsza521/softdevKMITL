package usecase

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	db_model "backend/internal/db_model"
	"backend/internal/notifications/dto"
	"backend/internal/notifications/interfaces"
)

type notificationUsecase struct {
	db   *gorm.DB
	repo interfaces.NotiRepository
}

func NewNotificationUsecase(db *gorm.DB, repo interfaces.NotiRepository) interfaces.NotiUsecase {
	return &notificationUsecase{db: db, repo: repo}
}

func (u *notificationUsecase) List(ctx context.Context, q dto.ListQuery) (*dto.ListResponse, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	sortAsc := strings.EqualFold(q.Sort, "created_at_asc")

	filter := interfaces.ListFilter{
		ReceiverID:   q.ReceiverID,
		ReceiverType: q.ReceiverType,
		IsRead:       q.IsRead,
		Offset:       (q.Page - 1) * q.PageSize,
		Limit:        q.PageSize,
		SortAsc:      sortAsc,
	}
	rows, total, err := u.repo.List(ctx, u.db, filter)
	if err != nil {
		return nil, err
	}

	items := make([]dto.Notification, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.Notification{
			ID:           r.ID,
			Title:        r.Title,
			Content:      r.Content,
			Type:         string(r.Type),
			ActionURL:    r.ActionURL,
			ReceiverID:   r.ReceiverID,
			ReceiverType: r.ReceiverType,
			IsRead:       r.IsRead,
			CreatedAt:    r.CreatedAt,
		})
	}

	totalPages := int((total + int64(q.PageSize) - 1) / int64(q.PageSize))
	return &dto.ListResponse{
		Items:      items,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (u *notificationUsecase) MarkRead(ctx context.Context, id uuid.UUID, isRead bool) error {
	return u.repo.MarkRead(ctx, u.db, id, isRead)
}

func (u *notificationUsecase) MarkAllRead(ctx context.Context, receiverId uuid.UUID, receiverType string) (int, error) {
	affected, err := u.repo.MarkAllRead(ctx, u.db, receiverId, receiverType)
	return int(affected), err
}

func (u *notificationUsecase) MockCreate(ctx context.Context, req dto.MockCreateRequest) (*dto.MockCreateResponse, error) {
	count := req.Count
	if count == 0 {
		count = 10
	}
	if count < 0 || count > 200 {
		return nil, errors.New("count must be between 1 and 200")
	}

	now := time.Now()
	samples := []struct {
		Title   string
		Content string
		Type    db_model.NotificationType
		Link    *string
	}{
		{"System Update", "ระบบจะปิดปรับปรุงเวลา 02:00-03:00 น.", db_model.NotificationTypeSystem, nil},
		{"New Booking", "คุณมีการจองใหม่ #BK-" + randSuffix(), db_model.NotificationTypeBooking, nil},
		{"Payment Received", "ชำระเงินสำเร็จสำหรับคำสั่งซื้อ #" + randSuffix(), db_model.NotificationTypePayment, nil},
	}

	notis := make([]db_model.Notifications, 0, count)
	for i := 0; i < count; i++ {
		s := samples[i%len(samples)]
		notis = append(notis, db_model.Notifications{
			Base:         db_model.Base{}, // ปล่อยให้ GORM เติม ID/CreatedAt เอง
			Title:        s.Title,
			Content:      s.Content,
			Type:         s.Type,
			ActionURL:    s.Link,
			ReceiverID:   req.ReceiverID,
			ReceiverType: req.ReceiverType,
			IsRead:       false,
			CreatedAt:    now.Add(-time.Duration(i) * time.Minute), // ไล่เวลาให้ดูสมจริง
		})
	}

	if err := u.repo.CreateBulk(ctx, u.db, notis); err != nil {
		return nil, err
	}
	return &dto.MockCreateResponse{Inserted: len(notis)}, nil
}

func randSuffix() string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
