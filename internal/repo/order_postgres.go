package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	v2 "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

const (
	// При отсутствии номера в базе добавляет запись, при конфликте - возвращает id пользователя
	createUploadedOrder = `
WITH _ AS (
    INSERT INTO upload_orders (user_id, number, status)
        VALUES ($1, $2, 'NEW')
        ON CONFLICT (number)
            DO NOTHING
        RETURNING user_id)
SELECT user_id
FROM upload_orders
WHERE number = $2;
`
	getUploadedOrdersByUser = `
SELECT number, COALESCE(status, '') AS status, COALESCE(accrual, 0) as accrual, created_at
    FROM upload_orders
WHERE user_id = $1
ORDER BY created_at DESC;
`
	updateUploadedOrder = `
UPDATE upload_orders
SET status = $2,
    accrual = $3
WHERE number = $1;
`
	createWithdrawnOrder = `
INSERT INTO withdraw_orders (user_id, number, sum)
VALUES ($1, $2, $3);
`
	getWithdrawOrdersByUser = `
SELECT number, sum, created_at
    FROM withdraw_orders
WHERE user_id = $1
ORDER BY created_at DESC;
`
)

var (
	stmtUpdateUploadedOrder *sqlx.Stmt
	stmtCreateWithdrawOrder *sqlx.Stmt
)

type OrderPostgres struct {
	db *v2.Postgre
}

func NewOrderPostgres(db *v2.Postgre) *OrderPostgres {
	var err error
	stmtUpdateUploadedOrder, err = db.Preparex(updateUploadedOrder)
	if err != nil {
		log.Printf("repo - NewOrderPostgres stmtUpdateUploadedOrder prepare: %v", err)
	}

	stmtUpdateCurrentBalance, err = db.Preparex(updateCurrentBalance)
	if err != nil {
		log.Printf("repo - NewOrderPostgres stmtUpdateCurrentBalance prepare: %v", err)
	}

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

	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) CreateUploadedOrder(userID int, orderNumber string) (int, error) {
	var userIDOut int
	row := r.db.QueryRow(createUploadedOrder, userID, orderNumber)
	if err := row.Scan(&userIDOut); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
	}

	return userIDOut, nil
}

func (r *OrderPostgres) GetUploadedOrders(ctx context.Context, userID int) ([]entity.UploadOrder, error) {
	var result []entity.UploadOrder

	ctxInner, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.db.SelectContext(ctxInner, &result, getUploadedOrdersByUser, userID); err != nil {
		return nil, fmt.Errorf("repo - get upload orders by user: %w", err)
	}

	return result, nil
}

func (r *OrderPostgres) UpdateUploadedOrder(number string, status string, accrual float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repo - start transaction: %w", err)
	}
	defer tx.Rollback()

	txStmtUpdateUploadedOrder := tx.StmtxContext(ctx, stmtUpdateUploadedOrder)
	txStmtUpdateBalance := tx.StmtxContext(ctx, stmtUpdateCurrentBalance)

	if _, err = txStmtUpdateUploadedOrder.Exec(number, status, accrual); err != nil {
		return fmt.Errorf("repo - txStmtUpdateUploadedOrder: %w", err)
	}

	if _, err = txStmtUpdateBalance.Exec(number, accrual); err != nil {
		return fmt.Errorf("repo - txStmtUpdateBalance: %w", err)
	}

	log.Println("repo - update upload order success")
	return tx.Commit()
}

func (r *OrderPostgres) CreateWithdrawOrder(userID int, orderNumber string, sum float32) error {
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

	before := balance
	if _, err = txStmtCreateWithdrawOrder.Exec(userID, orderNumber, sum); err != nil {
		return fmt.Errorf("repo - txStmtCreateWithdrawOrder: %w", err)
	}

	if _, err = txStmtUpdateWithdrawBalance.Exec(userID, sum); err != nil {
		return fmt.Errorf("repo - txStmtUpdateWithdrawBalance: %w", err)
	}

	if err = txStmtGetCurrentBalance.Get(&balance, userID); err != nil {
		return fmt.Errorf("repo - txStmtGetCurrentBalance: %w", err)
	}

	log.Println("repo - update withdraw order success: before=", before, " | after=", balance)
	return tx.Commit()
}

func (r *OrderPostgres) GetWithdrawOrders(ctx context.Context, userID int) ([]entity.WithdrawOrder, error) {
	var result []entity.WithdrawOrder

	ctxInner, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.db.SelectContext(ctxInner, &result, getWithdrawOrdersByUser, userID); err != nil {
		return nil, fmt.Errorf("repo - get upload orders by user: %w", err)
	}

	return result, nil

}
