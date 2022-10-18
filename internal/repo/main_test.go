package repo

import (
	"fmt"
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	postgres "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"os"
	"testing"
)

const (
	dbSource = "postgresql://admin:root@localhost:5432/postgres"
)

var testBalance *BalancePostgres

func TestMain(m *testing.M) {
	l := logger.New("repo_testing")

	pg, err := postgres.New(dbSource)
	if err != nil {
		l.Fatal(fmt.Errorf("cannot connect to db: %w", err))
	}
	defer pg.Close()

	testBalance = NewBalancePostgres(pg)

	os.Exit(m.Run())
}
