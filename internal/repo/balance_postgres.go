package repo

import v2 "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"

type BalancePostgres struct {
	db *v2.Postgre
}

func (b BalancePostgres) CreateBalance(userID int) error {
	//TODO implement me
	panic("implement me")
}

func (b BalancePostgres) UpdateBalance(userID int, sum float32) error {
	//TODO implement me
	panic("implement me")
}

func NewBalancePostgres(db *v2.Postgre) *BalancePostgres {
	return &BalancePostgres{db: db}
}
