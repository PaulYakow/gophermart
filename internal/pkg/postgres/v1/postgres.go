package v1

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

const (
	defaultMaxPoolSize     = 4
	defaultMaxConnIdleTime = time.Second * 30
	defaultMaxConnLifetime = time.Minute * 2

	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

//goland:noinspection SpellCheckingInspection
type Postgre struct {
	maxPollSize     int
	maxConnIdleTime time.Duration
	maxConnLifeTime time.Duration
	connAttempts    int
	connTimeout     time.Duration

	Pool *pgxpool.Pool
}

func New(dsn string, opts ...Option) (*Postgre, error) {
	pg := &Postgre{
		maxPollSize:     defaultMaxPoolSize,
		maxConnIdleTime: defaultMaxConnIdleTime,
		maxConnLifeTime: defaultMaxConnLifetime,
		connAttempts:    defaultConnAttempts,
		connTimeout:     defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgre - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPollSize)
	poolConfig.MaxConnIdleTime = pg.maxConnIdleTime
	poolConfig.MaxConnLifetime = pg.maxConnLifeTime

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgre - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgre) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
