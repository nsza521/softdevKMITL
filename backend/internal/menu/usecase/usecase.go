package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"backend/internal/db_model"
	"backend/internal/utils"
	iface "backend/internal/menu/interfaces"
	menu "backend/internal/menu/dto"
)

type menuUsecase struct{ 
	repo iface.MenuRepository 
	minioClient *minio.Client
}

func NewMenuUsecase(r iface.MenuRepository, minioClient *minio.Client) iface.MenuUsecase { 
	return &menuUsecase{
		repo: r, 
		minioClient: minioClient,
	} 
}

func (u *menuUsecase) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]iface.MenuItemBrief, error) {
	items, err := u.repo.ListMenuByRestaurant(ctx, restaurantID)
	if err != nil { return nil, err }
	out := make([]iface.MenuItemBrief, 0, len(items))
	for _, m := range items {
		typeIDs := make([]uuid.UUID, 0, len(m.MenuTypes))
		types := make([]iface.MenuTypeBrief, 0, len(m.MenuTypes))
		for _, t := range m.MenuTypes {
			typeIDs = append(typeIDs, t.ID)
			types = append(types, iface.MenuTypeBrief{ID: t.ID, Type: t.Type})
		}
		out = append(out, iface.MenuItemBrief{
			ID: m.ID, Name: m.Name, Price: m.Price, MenuPic: m.MenuPic,
			TimeTaken: m.TimeTaken, Description: m.Description,
			// MenuTypeIDs: typeIDs, 
			Types: types,
		})
	}
	return out, nil
}

func (u *menuUsecase) CheckRestaurantExists(ctx context.Context, restaurantID uuid.UUID) error {
	return u.repo.RestaurantExists(ctx, restaurantID)
}

// GetDetail retrieves detailed information about a menu item, including its associated menu types and add-ons.
// ใหม่: GetDetail
func (u *menuUsecase) GetDetail(itemID uuid.UUID) (menu.MenuItemDetailResponse, error) {
	item, err := u.repo.GetItemWithTypesAndAddOns(itemID)
	if err != nil {
		return menu.MenuItemDetailResponse{}, err
	}

	// 1) map types → []MenuTypeBrief
	types := make([]menu.MenuTypeBrief, 0, len(item.MenuTypes))
	for _, t := range item.MenuTypes {
		types = append(types, menu.MenuTypeBrief{ID: t.ID, Name: t.Type})
	}

	// 2) สร้างแผนที่ group โดย key = groupID (เพื่อ merge)
	type groupAgg struct {
		src    string               // "type" | "item" | "merged"
		group  models.MenuAddOnGroup
		opts   []models.MenuAddOnOption
	}
	groupMap := map[uuid.UUID]*groupAgg{}

	// 2.1 จาก type-level
	for _, mt := range item.MenuTypes {
		for _, g := range mt.AddOnGroups {
			entry := groupMap[g.ID]
			if entry == nil {
				cp := g // copy
				groupMap[g.ID] = &groupAgg{
					src:   "type",
					group: cp,
					opts:  append([]models.MenuAddOnOption{}, g.Options...),
				}
			} else {
				// merge options (กันซ้ำด้วย ID)
				entry.src = "merged"
				entry.opts = mergeOptions(entry.opts, g.Options)
			}
		}
	}

	// 2.2 จาก item-level (override/เพิ่ม)
	for _, g := range item.AddOnGroups {
		entry := groupMap[g.ID]
		if entry == nil {
			cp := g
			groupMap[g.ID] = &groupAgg{
				src:   "item",
				group: cp,
				opts:  append([]models.MenuAddOnOption{}, g.Options...),
			}
		} else {
			entry.src = "merged"
			// แนวนโยบายง่าย ๆ: รวม options เข้าด้วยกัน (dedupe ตาม ID)
			entry.opts = mergeOptions(entry.opts, g.Options)
			// ถ้าต้องการให้ item-level override ฟิลด์ group (เช่น Required/AllowQty) ก็อัปเดตที่นี่
			// ex: entry.group.Required = g.Required
		}
	}

	// 3) แปลงเป็น DTO
	addons := make([]menu.AddOnGroupDTO, 0, len(groupMap))
	for _, ag := range groupMap {
		optsDTO := make([]menu.AddOnOptionDTO, 0, len(ag.opts))
		for _, o := range ag.opts {
			optsDTO = append(optsDTO, menu.AddOnOptionDTO{
				ID:         o.ID,
				Name:       o.Name,
				PriceDelta: o.PriceDelta,
				IsDefault:  o.IsDefault,
				MaxQty:     o.MaxQty,
			})
		}
		addons = append(addons, menu.AddOnGroupDTO{
			ID:        ag.group.ID,
			Name:      ag.group.Name,
			Required:  ag.group.Required,
			MinSelect: ag.group.MinSelect,
			MaxSelect: ag.group.MaxSelect,
			AllowQty:  ag.group.AllowQty,
			From:      ag.src,
			Options:   optsDTO,
		})
	}

	// 4) สร้าง response
	resp := menu.MenuItemDetailResponse{
		ID:           item.ID,
		RestaurantID: item.RestaurantID,
		Name:         item.Name,
		Price:        item.Price,
		MenuPic:      item.MenuPic,
		TimeTaken:    item.TimeTaken,
		Description:  item.Description,
		Types:        types,
		AddOns:       addons,
	}
	return resp, nil
}

// helper: รวม option แบบกันซ้ำตาม ID
func mergeOptions(a, b []models.MenuAddOnOption) []models.MenuAddOnOption {
	idx := make(map[uuid.UUID]struct{}, len(a))
	out := make([]models.MenuAddOnOption, 0, len(a)+len(b))
	for _, x := range a {
		out = append(out, x)
		idx[x.ID] = struct{}{}
	}
	for _, x := range b {
		if _, ok := idx[x.ID]; !ok {
			out = append(out, x)
			idx[x.ID] = struct{}{}
		}
	}
	return out
}

func (u *menuUsecase) CreateMenuItem(ctx context.Context, restaurantID uuid.UUID, in *iface.CreateMenuItemRequest) (*iface.MenuItemBrief, error) {
	if in == nil { return nil, errors.New("nil request") }
	if err := u.repo.VerifyMenuTypesBelongToRestaurant(ctx, restaurantID, in.MenuTypeIDs); err != nil {
		return nil, err
	}
	mi := &models.MenuItem{
		Name: in.Name, Price: in.Price, 
		MenuPic: in.MenuPic,
		TimeTaken: in.TimeTaken, 
		Description: in.Description,
		RestaurantID: restaurantID,
	}
	if mi.TimeTaken == 0 { mi.TimeTaken = 1 }
	if err := u.repo.CreateMenuItem(ctx, mi); err != nil { return nil, err }
	if err := u.repo.AttachMenuTypes(ctx, mi.ID, in.MenuTypeIDs); err != nil { return nil, err }

	loaded, err := u.repo.LoadMenuItemWithTypes(ctx, mi.ID)
	if err != nil { return nil, err }
	resp := toBrief(loaded)
	return &resp, nil
}

func (u *menuUsecase) UpdateMenuItem(ctx context.Context, restaurantID uuid.UUID, id uuid.UUID, in *iface.UpdateMenuItemRequest) (*iface.MenuItemBrief, error) {
	if in == nil {
		return nil, errors.New("nil request")
	}

	fields := map[string]any{}
	if in.Name != nil       { fields["name"]        = *in.Name }
	if in.Price != nil      { fields["price"]       = *in.Price }
	if in.TimeTaken != nil  { fields["time_taken"]  = *in.TimeTaken }
	if in.Description != nil{ fields["description"] = *in.Description }
	if in.MenuPic != nil    { fields["menu_pic"]    = *in.MenuPic }

	if len(fields) > 0 {
		if err := u.repo.UpdateMenuItem(ctx, id, fields); err != nil {
			return nil, err
		}
	}

	// ถ้ามีการส่ง menu_type_ids → แทนที่ใหม่ทั้งหมด
	if in.MenuTypeIDs != nil {
		if err := u.repo.ReplaceMenuTypes(ctx, id, *in.MenuTypeIDs); err != nil {
			return nil, err
		}
	}

	// โหลดข้อมูลล่าสุดกลับมาให้ response
	loaded, err := u.repo.LoadMenuItemWithTypes(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toBrief(loaded)
	return &resp, nil
}


func (u *menuUsecase) DeleteMenuItem(ctx context.Context, restaurantID uuid.UUID, menuItemID uuid.UUID) error {
	return u.repo.DeleteMenuItem(ctx, menuItemID)
}

func toBrief(m *models.MenuItem) iface.MenuItemBrief {
	// typeIDs := make([]uuid.UUID, 0, len(m.MenuTypes))
	types := make([]iface.MenuTypeBrief, 0, len(m.MenuTypes))
	for _, t := range m.MenuTypes {
		// typeIDs = append(typeIDs, t.ID)
		types = append(types, iface.MenuTypeBrief{ID: t.ID, Type: t.Type})
	}
	return iface.MenuItemBrief{
		ID: m.ID, Name: m.Name, Price: m.Price, MenuPic: m.MenuPic,
		TimeTaken: m.TimeTaken, Description: m.Description,
		// MenuTypeIDs: typeIDs, 
		Types: types,
	}
}

func (u *menuUsecase) UploadMenuItemPicture(ctx context.Context, restaurantID uuid.UUID, itemID uuid.UUID, file *multipart.FileHeader) (string, error) {

	// Check if menu item exists
	menuItem, err := u.repo.GetMenuItemByID(ctx, itemID)
	if err != nil {
		return "", err
	}
	if menuItem.RestaurantID != restaurantID {
		return "", errors.New("menu item does not belong to this restaurant")
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	// Upload to MinIO
	const bucketName = "restaurant-pictures"
	const subBucket = "menu-items"
	filename := itemID.String()
	objectName := fmt.Sprintf("%s/%s", subBucket, filename)

	url, err := utils.UploadImage(fileContent, file, bucketName, objectName, u.minioClient)
	if err != nil {
		return "", err
	}

	// Update restaurant profile picture URL
	if menuItem != nil {
		menuItem.MenuPic = &url
	}
	err = u.repo.UpdateMenuItem(ctx, menuItem.ID, map[string]any{"menu_pic": menuItem.MenuPic})
	if err != nil {
		return "", err
	}

	// presignURL, err := utils.GetPresignedURL(u.minioClient, bucketName, objectName)
	// if err != nil {
	// 	return "", err
	// }

	return url, nil
}