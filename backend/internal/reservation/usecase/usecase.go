package usecase

import (

	"backend/internal/reservation/interfaces"
)

type TableReservationUsecase struct {
	tableReservationRepository interfaces.TableReservationRepository
}

func NewTableReservationUsecase(tableReservationRepository interfaces.TableReservationRepository) interfaces.TableReservationUsecase {
	return &TableReservationUsecase{
		tableReservationRepository: tableReservationRepository,
	}
}

