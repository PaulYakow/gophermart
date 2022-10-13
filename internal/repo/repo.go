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
    id          SERIAL PRIMARY KEY,
    user_id		INT NOT NULL,
    number      VARCHAR UNIQUE,
    status      VARCHAR,
    accrual     REAL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS balance
(
    id          SERIAL PRIMARY KEY,
    user_id		INT NOT NULL,
    current     REAL,
    withdrawn	REAL
);

CREATE TABLE IF NOT EXISTS withdraw_orders
(
    id          SERIAL PRIMARY KEY,
    user_id		INT NOT NULL,
    number      VARCHAR UNIQUE,
    sum			REAL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE upload_orders
    ADD FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE;

ALTER TABLE balance
    ADD FOREIGN KEY (user_id)
        REFERENCES users (id)
		ON DELETE CASCADE;

ALTER TABLE withdraw_orders
    ADD FOREIGN KEY (user_id)
        REFERENCES users (id)
		ON DELETE CASCADE;

CREATE INDEX ON upload_orders (user_id);
CREATE INDEX ON balance (user_id);
CREATE INDEX ON withdraw_orders (user_id);
`

type (
	IAuthorization interface {
		CreateUser(user entity.User) (int, error)
		GetUser(login, password string) (entity.User, error)
	}

	IUploadOrder interface {
		CreateUploadedOrder(userID int, orderNumber string) (int, error)
		GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrder, error)
		UpdateUploadedOrder(number string, status string, accrual float32) error
	}

	IBalance interface {
		GetBalance(ctx context.Context, userID int) (entity.Balance, error)
		UpdateCurrentBalance(userID int, sum float32) error
		UpdateWithdrawBalance(userID int, sum float32) error
	}

	// todo: add mutex
	Repo struct {
		IAuthorization
		IUploadOrder
		IBalance
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
		IUploadOrder:   NewOrderPostgres(db),
		IBalance:       NewBalancePostgres(db),
	}, nil
}
