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
    login         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS upload_orders
(
    number      VARCHAR UNIQUE,
    status      VARCHAR,
    accrual     NUMERIC,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id		INT REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS balance
(
    user_id		INT REFERENCES users (id) ON DELETE CASCADE,
    current     NUMERIC,
    withdrawn	NUMERIC
);

CREATE TABLE IF NOT EXISTS withdraw_orders
(
    user_id		INT REFERENCES users (id) ON DELETE CASCADE,
    number      VARCHAR UNIQUE,
    sum			NUMERIC,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON upload_orders (user_id);
CREATE INDEX ON balance (user_id);
CREATE INDEX ON withdraw_orders (user_id);
`

type (
	IAuthorization interface {
		CreateUser(login, passwordHash string) (int, error)
		GetUser(login string) (entity.UserDAO, error)
	}

	IOrder interface {
		CreateUploadedOrder(userID int, orderNumber string) (int, error)
		GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrderDAO, error)
		UpdateUploadedOrder(number string, status string, accrual float32) error

		CreateWithdrawOrder(userID int, orderNumber string, sum float32) error
		GetWithdrawOrders(ctx context.Context, userID int) ([]entity.WithdrawOrderDAO, error)
	}

	IBalance interface {
		GetBalance(ctx context.Context, userID int) (entity.BalanceDAO, error)
	}

	Repo struct {
		IAuthorization
		IOrder
		IBalance

		NotProcessedOrders []string
	}
)

func New(db *v2.Postgre) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("repo - New - create schema failed: %w", err)
	}

	authRepo, err := NewAuthPostgres(db)
	if err != nil {
		return nil, fmt.Errorf("repo - New - create repo/auth failed: %w", err)
	}

	orderRepo, err := NewOrderPostgres(db)
	if err != nil {
		return nil, fmt.Errorf("repo - New - create repo/order failed: %w", err)
	}

	return &Repo{
		IAuthorization:     authRepo,
		IOrder:             orderRepo,
		IBalance:           NewBalancePostgres(db),
		NotProcessedOrders: orderRepo.NotProcessedOrders,
	}, nil
}
