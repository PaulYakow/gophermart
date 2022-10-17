package service

import (
	"context"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	"github.com/PaulYakow/gophermart/internal/repo"
)

type (
	IAuthorization interface {
		CreateUser(user entity.User) (int, error)
		GenerateToken(login, password string) (string, error)
		ParseToken(token string) (int, error)
	}

	IOrder interface {
		CreateUploadedOrder(userID, orderNumber int) (int, error)
		GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrder, error)

		CreateWithdrawOrder(userID int, orderNumber string, sum float32) error
		GetWithdrawOrders(ctx context.Context, userID int) ([]entity.WithdrawOrder, error)
	}

	IBalance interface {
		GetBalance(ctx context.Context, userID int) (entity.Balance, error)
	}

	Service struct {
		IAuthorization
		IOrder
		IBalance

		Polling *PollService
	}
)

func NewService(repo *repo.Repo, logger logger.ILogger) *Service {
	return &Service{
		IAuthorization: NewAuthService(repo.IAuthorization),
		IOrder:         NewOrderService(repo.IOrder),
		IBalance:       NewBalanceService(repo.IBalance),
		Polling:        NewPollService(repo.IOrder, logger),
	}
}
