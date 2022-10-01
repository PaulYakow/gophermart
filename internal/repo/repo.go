package repo

import (
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/jmoiron/sqlx"
)

type (
	IAuthorization interface {
		CreateUser(user entity.User) (int, error)
		GetUser(login, password string) (entity.User, error)
	}

	IOrder interface {
		Create(userId, orderNumber int) (int, error)
	}

	Repo struct {
		IAuthorization
		IOrder
	}
)

func New(db *sqlx.DB) *Repo {
	return &Repo{
		IAuthorization: NewAuthPostgres(db),
	}
}
