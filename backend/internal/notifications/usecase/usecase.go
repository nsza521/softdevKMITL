package usecase

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"
	// "fmt"
	"encoding/json"

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
        notification := dto.Notification{
            ID:           r.ID,
            Title:        r.Title,
            Content:      r.Content,
            Type:         string(r.Type),
            ActionURL:    r.ActionURL,
            ReceiverID:   r.ReceiverID,
            ReceiverType: r.ReceiverType,
            IsRead:       r.IsRead,
            CreatedAt:    r.CreatedAt.Format("02-01-2006 15:04"),
        }

        // Parse attributes จาก JSON string
        if r.Attributes != nil && *r.Attributes != "" {
            var attrs map[string]interface{}
            if err := json.Unmarshal([]byte(*r.Attributes), &attrs); err == nil {
                notification.Attributes = attrs
            }
        }

        items = append(items, notification)
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

func (u *notificationUsecase) CreateFromEvent(ctx context.Context, req dto.CreateEventRequest) (*dto.CreateEventResponse, error) {
	
	if req.ReceiverID == uuid.Nil && req.ReceiverUsername != "" {
		switch strings.ToLower(req.ReceiverType) {
		case "customer":
			var cust db_model.Customer
			if err := u.db.WithContext(ctx).Where("username = ?", req.ReceiverUsername).First(&cust).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("customer username not found")
				}
				return nil, err
			}
			req.ReceiverID = cust.ID
		case "restaurant":
            var rest db_model.Restaurant
            if err := u.db.WithContext(ctx).Where("username = ?", req.ReceiverUsername).First(&rest).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                    return nil, errors.New("restaurant username not found")
                }
                return nil, err
            }
            req.ReceiverID = rest.ID
        default:
            return nil, errors.New("invalid receiverType")
		}
	}
	
	var title, content, actionURL string
	var attributes map[string]interface{}

	switch req.Event {
	case "reserve_with":
        d := req.Data.(map[string]interface{})
        title = "คุณได้รับคำเชิญจาก " + firstString(d["members"])
        content = "คุณได้รับคำเชิญให้เข้าร่วมการจองโต๊ะ"
        
        // แยกข้อมูลเป็น attributes
        attributes = map[string]interface{}{
            "tableNo":    d["tableNo"],
            "when":       d["when"],
            "restaurant": d["restaurant"],
            "members":    toStrings(d["members"]),
            "inviter":    firstString(d["members"]),
        }

    case "order_finished":
        d := req.Data.(map[string]interface{})
        title = "อาหารพร้อมแล้ว !"
        content = "คุณสามารถรับอาหารได้ที่ร้าน"
        
        attributes = map[string]interface{}{
            "tableNo":    d["tableNo"],
            "when":       d["when"],
            "restaurant": d["restaurant"],
			"menu": 	  d["menu"],
            "queueNo":    d["queueNo"],
        }

    case "order_canceled":
        d := req.Data.(map[string]interface{})
        title = "ออเดอร์ถูกยกเลิก"
        content = "คำสั่งซื้อของคุณถูกยกเลิก"
        
        attributes = map[string]interface{}{
            "tableNo":    d["tableNo"],
            "when":       d["when"],
            "restaurant": d["restaurant"],
            "reason":     d["reason"],
        }

    case "reserve_success":
        d := req.Data.(map[string]interface{})
        title = "จองโต๊ะสำเร็จ !"
        content = "การจองโต๊ะของคุณสำเร็จแล้ว"
        
        attributes = map[string]interface{}{
            "tableNo":    d["tableNo"],
            "when":       d["when"],
            "restaurant": d["restaurant"],
            "seat":       d["seat"],
        }

    case "reserve_failed":
        d := req.Data.(map[string]interface{})
        title = "จองโต๊ะไม่สำเร็จ"
        content = "ไม่สามารถจองโต๊ะได้ในขณะนี้"
        
        attributes = map[string]interface{}{
            "tableNo":    d["tableNo"],
            "when":       d["when"],
            "restaurant": d["restaurant"],
        }

	default:
		return nil, errors.New("unknown event")
	}

	// Convert attributes to JSON string for database storage
    var attributesJSON *string
    if attributes != nil {

		if whenVal, ok := attributes["when"]; ok {
			attributes["when"] = formatWhen(whenVal)
		}

        if jsonBytes, err := json.Marshal(attributes); err == nil {
            jsonStr := string(jsonBytes)
            attributesJSON = &jsonStr
        }
    }

	noti := db_model.Notifications{
		Title:        title,
		Content:      content,
		Type:         db_model.NotificationType(strings.ToUpper(req.Event)), // หรือ map เป็น enum ของคุณ
		ReceiverID:   req.ReceiverID,
		ReceiverType: req.ReceiverType,
		IsRead:       false,
		ActionURL:    strPtrOrNil(actionURL),
		Attributes: attributesJSON,
		CreatedAt:    time.Now(),
	}

	if err := u.repo.Create(ctx, u.db, &noti); err != nil {
		return nil, err
	}

	return &dto.CreateEventResponse{
		ID:        noti.ID,
		Title:     noti.Title,
		Content:   noti.Content,
		Attributes: attributes,
		CreatedAt: noti.CreatedAt.Format("02-01-2006 15:04"),
	}, nil
}

// helpers
func firstString(v interface{}) string {
	if a, ok := v.([]interface{}); ok && len(a) > 0 {
		if s, ok := a[0].(string); ok { return s }
	}
	if s, ok := v.(string); ok { return s }
	return ""
}
func toStrings(v interface{}) []string {
	var out []string
	if a, ok := v.([]interface{}); ok {
		for _, x := range a { if s, ok := x.(string); ok { out = append(out, s) } }
	}
	return out
}
func strPtrOrNil(s string) *string { if s == "" { return nil }; return &s }

// formatWhen พยายาม parse ค่าวันเวลาแล้วคืนเป็น "02-01-2006 15:04"
// ถา้ parse ไม่ได้ จะคืนค่าเดิม (string)
func formatWhen(v interface{}) string {
    // รองรับ time.Time
    if t, ok := v.(time.Time); ok {
        return t.Format("02-01-2006 15:04")
    }
    // รองรับ string
    s, ok := v.(string)
    if !ok {
        // ถ้าไม่ใช่ string ให้คืนค่าที่เป็น string โดยตรง
        return ""
    }

    layouts := []string{
        time.RFC3339,
        "2006-01-02 15:04:05",
        "2006-01-02 15:04",
        "02-01-2006 15:04",
        "02/01/2006 15:04",
        "02 Jan 2006 15:04",
        "02 Jan 2006, 15:04",
        "02 Jan 2006",
        "02 Jan 2006 15:04:05",
    }

    for _, l := range layouts {
        if tt, err := time.Parse(l, s); err == nil {
            return tt.Format("02-01-2006 15:04")
        }
    }

    // ถ้า parse ไม่ได้ คืน string เดิม (fallback)
    return s
}