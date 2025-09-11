package usecase

import (

	"backend/internal/table/interfaces"
)

type TableUsecase struct {
	tableRepository interfaces.TableRepository
}

func NewTableUsecase(tableRepository interfaces.TableRepository) interfaces.TableUsecase {
	return &TableUsecase{
		tableRepository: tableRepository,
	}
}

