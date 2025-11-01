package usecase

import (
	"backend/internal/menu/dto"
	"backend/internal/menu/interfaces"
	 models "backend/internal/db_model"
	"github.com/google/uuid"
)

type addOnUsecase struct {
	repo interfaces.AddOnRepository
}

func NewAddOnUsecase(r interfaces.AddOnRepository) interfaces.AddOnUsecase {
	return &addOnUsecase{repo: r}
}

// Group
func (u *addOnUsecase) CreateGroup(restaurantID uuid.UUID, input menu.CreateAddOnGroupRequest) (menu.AddOnGroupResponse, error) {
	group := models.MenuAddOnGroup{
		RestaurantID: restaurantID,
		Name:         input.Name,
		Required:     input.Required,
		MinSelect:    input.MinSelect,
		MaxSelect:    input.MaxSelect,
		AllowQty:     input.AllowQty,
	}
	if err := u.repo.CreateGroup(&group); err != nil {
		return menu.AddOnGroupResponse{}, err
	}
	return menu.AddOnGroupResponse{
		ID:         group.ID,
		Name:       group.Name,
		Required:   group.Required,
		MinSelect:  group.MinSelect,
		MaxSelect:  group.MaxSelect,
		AllowQty:   group.AllowQty,
		Restaurant: group.RestaurantID,
	}, nil
}

func (u *addOnUsecase) ListGroups(restaurantID uuid.UUID) ([]menu.AddOnGroupResponse, error) {
	groups, err := u.repo.GetGroupsByRestaurant(restaurantID)
	if err != nil {
		return nil, err
	}
	var res []menu.AddOnGroupResponse
	for _, g := range groups {
		res = append(res, menu.AddOnGroupResponse{
			ID:         g.ID,
			Name:       g.Name,
			Required:   g.Required,
			MinSelect:  g.MinSelect,
			MaxSelect:  g.MaxSelect,
			AllowQty:   g.AllowQty,
			Restaurant: g.RestaurantID,
		})
	}
	return res, nil
}

func (u *addOnUsecase) UpdateGroup(id uuid.UUID, input menu.UpdateAddOnGroupRequest) error {
	group, err := u.repo.GetGroupByID(id)
	if err != nil {
		return err
	}
	if input.Name != nil {
		group.Name = *input.Name
	}
	if input.Required != nil {
		group.Required = *input.Required
	}
	if input.MinSelect != nil {
		group.MinSelect = input.MinSelect
	}
	if input.MaxSelect != nil {
		group.MaxSelect = input.MaxSelect
	}
	if input.AllowQty != nil {
		group.AllowQty = *input.AllowQty
	}
	return u.repo.UpdateGroup(group)
}

func (u *addOnUsecase) DeleteGroup(id uuid.UUID) error {
	return u.repo.DeleteGroup(id)
}



func (u *addOnUsecase) LinkGroupToTypes(groupID uuid.UUID, typeIDs []uuid.UUID) error {
    return u.repo.LinkGroupToTypes(groupID, typeIDs)
}
func (u *addOnUsecase) UnlinkGroupFromType(groupID, typeID uuid.UUID) error {
    return u.repo.UnlinkGroupFromType(groupID, typeID)
}

// Option

// usecase/addon_uc.go
func (u *addOnUsecase) GetOption(id uuid.UUID) (menu.AddOnOptionResponse, error) {
    opt, err := u.repo.GetOptionByID(id)
    if err != nil {
        return menu.AddOnOptionResponse{}, err
    }
    return menu.AddOnOptionResponse{
        ID: opt.ID, Name: opt.Name, PriceDelta: opt.PriceDelta,
        IsDefault: opt.IsDefault, MaxQty: opt.MaxQty, GroupID: opt.GroupID,
    }, nil
}

func (u *addOnUsecase) ListOptions(groupID uuid.UUID) ([]menu.AddOnOptionResponse, error) {
    group, err := u.repo.GetGroupByID(groupID)
    if err != nil {
        return nil, err
    }
    var res []menu.AddOnOptionResponse
    for _, opt := range group.Options {
        res = append(res, menu.AddOnOptionResponse{
            ID: opt.ID, Name: opt.Name, PriceDelta: opt.PriceDelta,
            IsDefault: opt.IsDefault, MaxQty: opt.MaxQty, GroupID: opt.GroupID,
        })
    }
    return res, nil
}


func (u *addOnUsecase) CreateOption(groupID uuid.UUID, input menu.CreateAddOnOptionRequest) (menu.AddOnOptionResponse, error) {
	opt := models.MenuAddOnOption{
		GroupID:    groupID,
		Name:       input.Name,
		PriceDelta: input.PriceDelta,
		IsDefault:  input.IsDefault,
		MaxQty:     input.MaxQty,
	}
	if err := u.repo.CreateOption(&opt); err != nil {
		return menu.AddOnOptionResponse{}, err
	}
	return menu.AddOnOptionResponse{
		ID:         opt.ID,
		Name:       opt.Name,
		PriceDelta: opt.PriceDelta,
		IsDefault:  opt.IsDefault,
		MaxQty:     opt.MaxQty,
		GroupID:    opt.GroupID,
	}, nil
}

func (u *addOnUsecase) UpdateOption(id uuid.UUID, input menu.UpdateAddOnOptionRequest) error {
	opt, err := u.repo.GetOptionByID(id)
	if err != nil {
		return err
	}
	if input.Name != nil {
		opt.Name = *input.Name
	}
	if input.PriceDelta != nil {
		opt.PriceDelta = *input.PriceDelta
	}
	if input.IsDefault != nil {
		opt.IsDefault = *input.IsDefault
	}
	if input.MaxQty != nil {
		opt.MaxQty = input.MaxQty
	}
	return u.repo.UpdateOption(opt)
}

func (u *addOnUsecase) DeleteOption(id uuid.UUID) error {
	return u.repo.DeleteOption(id)
}
