package usecase

import (

	"backend/internal/menu/interfaces"
)

type MenuUsecase struct {
	menuRepository interfaces.MenuRepository
}

func NewMenuUsecase(menuRepository interfaces.MenuRepository) interfaces.MenuUsecase {
	return &MenuUsecase{
		menuRepository: menuRepository,
	}
}

