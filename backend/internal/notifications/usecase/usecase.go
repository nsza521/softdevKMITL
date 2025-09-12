package usecase

import (

	"backend/internal/notifications/interfaces"
)

type NotiUsecase struct {
	notiRepository interfaces.NotiRepository
}

func NewNotiUsecase(notiRepository interfaces.NotiRepository) interfaces.NotiUsecase {
	return &NotiUsecase{
		notiRepository: notiRepository,
	}
}

