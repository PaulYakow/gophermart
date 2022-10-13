package repo

import (
	"context"
	"github.com/PaulYakow/gophermart/internal/entity"
	v2 "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
)

const (
	createBalanceForUser = `
INSERT INTO balance (user_id, current, withdrawn)
VALUES ($1, 0, 0);
`
	getCurrentBalanceByUser = `
SELECT current
FROM balance
WHERE user_id = $1;
`
	updateBalance = `
UPDATE balance
SET current = current + $2
WHERE user_id = (SELECT user_id FROM upload_orders WHERE number = $1);
`
)

type BalancePostgres struct {
	db *v2.Postgre
}

func NewBalancePostgres(db *v2.Postgre) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func (r BalancePostgres) GetBalance(ctx context.Context, userID int) (entity.Balance, error) {
	var balance entity.Balance
	err := r.db.Get(&balance, getCurrentBalanceByUser, userID)
	return balance, err
}

func (r BalancePostgres) UpdateCurrentBalance(userID int, sum float32) error {
	//TODO implement me
	panic("implement me")
}

func (r BalancePostgres) UpdateWithdrawBalance(userID int, sum float32) error {
	//TODO implement me
	panic("implement me")
}
