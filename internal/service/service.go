package service

import (
	"context"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	"github.com/PaulYakow/gophermart/internal/repo"
)

type (
	IAuthorization interface {
		CreateUser(login, password string) (int, error)
		GetUser(login, password string) (int, error)
		GenerateToken(userID int) (string, error)
		ParseToken(token string) (int, error)
	}

	IOrder interface {
		CreateUploadedOrder(userID, orderNumber int) (int, error)
		GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrderDTO, error)

		CreateWithdrawOrder(userID int, orderNumber string, sum float32) error
		GetWithdrawOrders(ctx context.Context, userID int) ([]entity.WithdrawOrderDTO, error)
	}

	IBalance interface {
		GetBalance(ctx context.Context, userID int) (entity.BalanceDTO, error)
	}

	Service struct {
		IAuthorization
		IOrder
		IBalance

		Polling *PollService
	}
)

func NewService(repo *repo.Repo, logger logger.ILogger) (*Service, error) {
	authService, err := NewAuthService(repo.IAuthorization)
	if err != nil {
		return nil, fmt.Errorf("service - New - create auth failed: %w", err)
	}

	return &Service{
		IAuthorization: authService,
		IOrder:         NewOrderService(repo.IOrder),
		IBalance:       NewBalanceService(repo.IBalance),
		Polling:        NewPollService(repo.IOrder, logger),
	}, nil
}
