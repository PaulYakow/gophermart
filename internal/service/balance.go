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

func (s BalanceService) GetBalance(ctx context.Context, userID int) (entity.Balance, error) {
	return s.repo.GetBalance(ctx, userID)
}

func (s BalanceService) UpdateCurrentBalance(userID int, sum float32) error {
	return s.repo.UpdateCurrentBalance(userID, sum)
}

func (s BalanceService) UpdateWithdrawBalance(userID int, sum float32) error {
	return s.repo.UpdateWithdrawBalance(userID, sum)
}
