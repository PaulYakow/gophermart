package service

import (
	"context"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/repo"
)

type BalanceService struct {
	repo repo.IBalance
}

func NewBalanceService(repo repo.IBalance) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID int) (entity.BalanceDTO, error) {
	balance, err := s.repo.GetBalance(ctx, userID)
	if err != nil {
		return entity.BalanceDTO{}, err
	}

	return entity.BalanceDTO{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}, nil
}
