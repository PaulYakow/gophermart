package repo

import (
	"context"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"github.com/jmoiron/sqlx"
)

const (
	createBalanceForUser = `
INSERT INTO balance (user_id, current, withdrawn)
VALUES ($1, 0, 0);
`
	getBalanceByUser = `
SELECT current, withdrawn
FROM balance
WHERE user_id = $1;
`
	updateCurrentBalance = `
UPDATE balance
SET current = current + $2
WHERE user_id = (SELECT user_id FROM upload_orders WHERE number = $1);
`
	updateWithdrawBalance = `
UPDATE balance
SET current = current - $2,
    withdrawn = withdrawn + $2
WHERE user_id = $1;
`
)

var (
	stmtGetBalance            *sqlx.Stmt
	stmtUpdateCurrentBalance  *sqlx.Stmt
	stmtUpdateWithdrawBalance *sqlx.Stmt
)

type BalancePostgres struct {
	db *v2.Postgre
}

func NewBalancePostgres(db *v2.Postgre) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func (r *BalancePostgres) GetBalance(ctx context.Context, userID int) (entity.Balance, error) {
	var balance entity.Balance
	err := r.db.GetContext(ctx, &balance, getBalanceByUser, userID)
	return balance, err
}
