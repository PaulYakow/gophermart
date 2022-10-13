package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
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
	stmtGetCurrentBalance     *sqlx.Stmt
	stmtUpdateCurrentBalance  *sqlx.Stmt
	stmtUpdateWithdrawBalance *sqlx.Stmt
)

type BalancePostgres struct {
	db *v2.Postgre
}

func NewBalancePostgres(db *v2.Postgre) *BalancePostgres {
	var err error
	stmtGetCurrentBalance, err = db.Preparex(getCurrentBalanceByUser)
	if err != nil {
		log.Printf("repo - NewOrderPostgres stmtGetCurrentBalance prepare: %v", err)
	}

	stmtCreateWithdrawOrder, err = db.Preparex(createWithdrawnOrder)
	if err != nil {
		log.Printf("repo - NewOrderPostgres stmtCreateWithdrawOrder prepare: %v", err)
	}

	stmtUpdateWithdrawBalance, err = db.Preparex(updateWithdrawBalance)
	if err != nil {
		log.Printf("repo - NewOrderPostgres stmtUpdateWithdrawBalance prepare: %v", err)
	}

	return &BalancePostgres{db: db}
}

func (r BalancePostgres) GetBalance(ctx context.Context, userID int) (entity.Balance, error) {
	var balance entity.Balance
	err := r.db.Get(&balance, getCurrentBalanceByUser, userID)
	return balance, err
}

func (r BalancePostgres) UpdateWithdrawBalance(userID, orderNumber int, sum float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repo - start transaction: %w", err)
	}
	defer tx.Rollback()

	txStmtGetCurrentBalance := tx.StmtxContext(ctx, stmtGetCurrentBalance)
	txStmtCreateWithdrawOrder := tx.StmtxContext(ctx, stmtCreateWithdrawOrder)
	txStmtUpdateWithdrawBalance := tx.StmtxContext(ctx, stmtUpdateWithdrawBalance)

	var balance entity.Balance
	if err = txStmtGetCurrentBalance.Get(&balance, userID); err != nil {
		return fmt.Errorf("repo - txStmtGetCurrentBalance: %w", err)
	}

	if balance.Current-sum < 0 {
		return ErrNoFunds
	}

	if _, err = txStmtCreateWithdrawOrder.Exec(userID, orderNumber, sum); err != nil {
		return fmt.Errorf("repo - txStmtCreateWithdrawOrder: %w", err)
	}

	if _, err = txStmtUpdateWithdrawBalance.Exec(userID, sum); err != nil {
		return fmt.Errorf("repo - txStmtUpdateWithdrawBalance: %w", err)
	}

	log.Println("repo - update withdraw order success")
	return tx.Commit()
}
