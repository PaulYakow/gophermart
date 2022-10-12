package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	v2 "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
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
	getUploadedOrderByUser = `
SELECT number, status, accrual, created_at
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
)

type OrderPostgres struct {
	db *v2.Postgre
}

func NewOrderPostgres(db *v2.Postgre) *OrderPostgres {
	// todo: Named stmt

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

	if err := r.db.SelectContext(ctxInner, &result, getUploadedOrderByUser, userID); err != nil {
		return nil, fmt.Errorf("repo - get upload orders by user: %w", err)
	}

	return result, nil
}

func (r *OrderPostgres) UpdateUploadedOrder(number string, status string, accrual float32) error {
	result, err := r.db.Exec(updateUploadedOrder, number, status, accrual)
	if err != nil {
		return fmt.Errorf("repo - update upload order: %w", err)
	}

	fmt.Println("repo - update upload order success: ", result)
	return nil
}
