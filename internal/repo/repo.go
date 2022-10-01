package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"time"
)

type (
	IAuthorization interface {
		CreateUser(user entity.User) (int, error)
		GetUser(login, password string) (entity.User, error)
	}

	IOrder interface {
		Create(userID, orderNumber int) (int, error)
	}

	Repo struct {
		IAuthorization
		IOrder
	}
)

func New(db *v2.Postgre) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, schemaUsers)
	if err != nil {
		return nil, fmt.Errorf("repo - New - create table failed: %w", err)
	}

	return &Repo{
		IAuthorization: NewAuthPostgres(db),
	}, nil
}
