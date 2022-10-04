package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"time"
)

const schema = `
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    login         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS upload_orders
(
    id          SERIAL PRIMARY KEY,
    user_id		SERIAL,
    number      BIGINT UNIQUE,
    status      VARCHAR(255),
    accrual     REAL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS balance
(
    id          SERIAL PRIMARY KEY,
    user_id		SERIAL,
    current     REAL,
    withdrawn	REAL
);

CREATE TABLE IF NOT EXISTS withdraw_orders
(
    id          SERIAL PRIMARY KEY,
    user_id		SERIAL,
    number      BIGINT UNIQUE,
    sum			REAL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

type (
	IAuthorization interface {
		CreateUser(user entity.User) (int, error)
		GetUser(login, password string) (entity.User, error)
	}

	IOrder interface {
		CreateOrder(userID, orderNumber int) (int, error)
	}

	Repo struct {
		IAuthorization
		IOrder
	}
)

func New(db *v2.Postgre) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - New - create schema failed: %w", err)
	}

	return &Repo{
		IAuthorization: NewAuthPostgres(db),
		IOrder:         NewOrderPostgres(db),
	}, nil
}
