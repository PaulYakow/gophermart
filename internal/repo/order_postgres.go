package repo

import (
	"database/sql"
	"errors"
	v2 "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
)

const (
	// При отсутствии номера в базе добавляет запись, при конфликте - возвращает id пользователя
	createOrder = `
WITH _ AS (
    INSERT INTO upload_orders (user_id, number)
        VALUES ($1, $2)
        ON CONFLICT (number)
            DO NOTHING
        RETURNING user_id)
SELECT user_id
FROM upload_orders
WHERE number = $2;
`
)

type OrderPostgres struct {
	db *v2.Postgre
}

func NewOrderPostgres(db *v2.Postgre) *OrderPostgres {
	// todo: Named stmt

	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) CreateOrder(userID, orderNumber int) (int, error) {
	var userIDOut int
	row := r.db.QueryRow(createOrder, userID, orderNumber)
	if err := row.Scan(&userIDOut); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
	}

	return userIDOut, nil
}
