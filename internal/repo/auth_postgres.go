package repo

import (
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
)

// todo: Named stmt

const (
	createUser = `
INSERT INTO users (login, password_hash)
VALUES ($1, $2)
RETURNING id;
`
	getUser = `
SELECT id
FROM users
WHERE login=$1 AND password_hash=$2;
`
)

type AuthPostgres struct {
	db *v2.Postgre
}

func NewAuthPostgres(db *v2.Postgre) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	var id int
	row := r.db.QueryRow(createUser, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, getUser, login, password)
	return user, err
}
