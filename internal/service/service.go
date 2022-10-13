package service

import (
	"context"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/repo"
)

type (
	IAuthorization interface {
		CreateUser(user entity.User) (int, error)
		GenerateToken(login, password string) (string, error)
		ParseToken(token string) (int, error)
	}

	IUploadOrder interface {
		CreateUploadedOrder(userID, orderNumber int) (int, error)
		GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrder, error)
	}

	IBalance interface {
		GetBalance(ctx context.Context, userID int) (entity.Balance, error)
		UpdateCurrentBalance(userID int, sum float32) error
		UpdateWithdrawBalance(userID int, sum float32) error
	}

	Service struct {
		IAuthorization
		IUploadOrder
		IBalance

		Polling *PollService
	}
)

func NewService(repo *repo.Repo, pollingAddress string) *Service {
	// todo: нет возможности вести логи - пробросить сюда логгер
	// todo: при рестарте сервиса реализовать перезапуск опроса тех заказов, статус которых не окончательный
	return &Service{
		IAuthorization: NewAuthService(repo.IAuthorization),
		IUploadOrder:   NewOrderService(repo.IUploadOrder),
		IBalance:       NewBalanceService(repo.IBalance),
		Polling:        NewPollService(repo.IUploadOrder, pollingAddress),
	}
}
