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

var (
	stmtCreateUser  *sqlx.Stmt
	stmtInitBalance *sqlx.Stmt
)

func NewAuthPostgres(db *v2.Postgre) *AuthPostgres {
	var err error
	stmtCreateUser, err = db.Preparex(createUser)
	if err != nil {
		log.Printf("repo - NewAuthPostgres stmtCreateUser prepare: %v", err)
	}

	stmtInitBalance, err = db.Preparex(createBalanceForUser)
	if err != nil {
		log.Printf("repo - NewAuthPostgres stmtInitBalance prepare: %v", err)
	}

	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("repo - start transaction: %w", err)
	}
	defer tx.Rollback()

	txStmtCreateUser := tx.StmtxContext(ctx, stmtCreateUser)
	txStmtInitBalance := tx.StmtxContext(ctx, stmtInitBalance)

	var id int
	if err = txStmtCreateUser.Get(&id, user.Login, user.Password); err != nil {
		return 0, fmt.Errorf("repo - txStmtCreateUser: %w", err)
	}

	if _, err = txStmtInitBalance.Exec(id); err != nil {
		return 0, fmt.Errorf("repo - txStmtInitBalance: %w", err)
	}

	return id, tx.Commit()
}

func (r *AuthPostgres) GetUser(login, password string) (entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, getUser, login, password)
	return user, err
}
