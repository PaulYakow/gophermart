package service

import (
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/repo"
)

type (
	IAuthorization interface {
		CreateUser(user entity.User) (int, error)
		GenerateToken(login, password string) (string, error)
		ParseToken(token string) (int, error)
	}

	IOrder interface {
		CreateOrder(userID, orderNumber int) (int, error)
	}

	Service struct {
		IAuthorization
		IOrder
	}
)

func NewService(repo *repo.Repo) *Service {
	return &Service{
		IAuthorization: NewAuthService(repo.IAuthorization),
		IOrder:         NewOrderService(repo.IOrder),
	}
}
