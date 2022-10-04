package repo

import v2 "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"

const (
	createOrder = `
INSERT INTO upload_orders (user_id, number)
VALUES ($1, $2)
RETURNING id;
`
	getOrderByUser = `

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
	var orderID int
	row := r.db.QueryRow(createOrder, userID, orderNumber)
	if err := row.Scan(&orderID); err != nil {
		return 0, err
	}

	return orderID, nil
}
