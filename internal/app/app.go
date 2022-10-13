package app

import (
	"context"
	"fmt"
	"github.com/PaulYakow/gophermart/config"
	ctrl "github.com/PaulYakow/gophermart/internal/controller/v1"
	"github.com/PaulYakow/gophermart/internal/pkg/httpserver"
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	postgres "github.com/PaulYakow/gophermart/internal/pkg/postgres/v2"
	"github.com/PaulYakow/gophermart/internal/repo"
	"github.com/PaulYakow/gophermart/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Cfg) {
	l := logger.New()
	defer l.Exit()

	// Postgres storage
	pg, err := postgres.New(cfg.Dsn)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - failed to initialize db: %w", err))
	}
	defer pg.Close()

	repos, err := repo.New(pg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - repo.New: %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Services
	// todo: слишком жёстко привязан конечный маршрут - а если потребуется несколько эндпойнтов опрашивать?
	endpoint := "http://" + cfg.AccrualAddress + "/api/orders"
	services := service.NewService(repos, endpoint)
	go services.Polling.Run(ctx)

	//HTTP server
	handler := ctrl.NewHandler(services, l)
	srv := httpserver.New(handler.InitRoutes(), httpserver.Address(cfg.Address))

	l.Info("app - Run - params: a=%s | d=%s | r=%s",
		cfg.Address, cfg.Dsn, cfg.AccrualAddress)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %v", s.String())
	case err := <-srv.Notify():
		l.Error(fmt.Errorf("app - Run - Notify: %w", err))
	}

	// Shutdown
	cancel()
	err = srv.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - Shutdown: %w", err))
	}
}
