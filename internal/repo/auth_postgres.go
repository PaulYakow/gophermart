package repo

import (
	"context"
	"fmt"
	"github.com/PaulYakow/gophermart/internal/entity"
	"github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"time"
)

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

// todo: add mutex
type AuthPostgres struct {
	db *v2.Postgre
}

var (
	stmtCreateUser  *sqlx.Stmt
	stmtInitBalance *sqlx.Stmt
)

func NewAuthPostgres(db *v2.Postgre) (*AuthPostgres, error) {
	var err error
	stmtCreateUser, err = db.Preparex(createUser)
	if err != nil {
		return nil, fmt.Errorf("repo/auth - NewAuthPostgres stmtCreateUser prepare: %v", err)
	}

	stmtInitBalance, err = db.Preparex(createBalanceForUser)
	if err != nil {
		return nil, fmt.Errorf("repo/auth - NewAuthPostgres stmtInitBalance prepare: %v", err)
	}

	return &AuthPostgres{db: db}, nil
}

func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("repo/auth - start transaction: %w", err)
	}
	defer tx.Rollback()

	txStmtCreateUser := tx.StmtxContext(ctx, stmtCreateUser)
	txStmtInitBalance := tx.StmtxContext(ctx, stmtInitBalance)

	var id int
	if err = txStmtCreateUser.Get(&id, user.Login, user.Password); err != nil {
		pqErr := err.(pgx.PgError)
		if pqErr.Code == "23505" {
			return 0, ErrDuplicateKey
		}
		return 0, fmt.Errorf("repo/auth - txStmtCreateUser: %w", err)
	}

	if _, err = txStmtInitBalance.Exec(id); err != nil {
		return 0, fmt.Errorf("repo/auth - txStmtInitBalance: %w", err)
	}

	return id, tx.Commit()
}

func (r *AuthPostgres) GetUser(login, password string) (entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, getUser, login, password)
	return user, err
}
