package usecase

import "github.com/osniantonio/challenge-clean-architecture/internal/entity"

type ListOrdersOutputDTO []entity.Order

type ListOrdersUseCase struct {
	repository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(repository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{repository}
}

func (uc *ListOrdersUseCase) Execute() (ListOrdersOutputDTO, error) {
	orders, err := uc.repository.FindAll()
	if err != nil {
		return ListOrdersOutputDTO{}, err
	}
	return orders, nil
}
