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

	Service struct {
		IAuthorization
		IUploadOrder
	}
)

func NewService(repo *repo.Repo) *Service {
	return &Service{
		IAuthorization: NewAuthService(repo.IAuthorization),
		IUploadOrder:   NewOrderService(repo.IUploadOrder),
	}
}
