package app

import (
	"context"
	"fmt"
	"github.com/PaulYakow/gophermart/config"
	"github.com/PaulYakow/gophermart/internal/controller"
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
	appLogger := logger.New("app")
	defer appLogger.Exit()

	// Postgres storage
	pg, err := postgres.New(cfg.Dsn)
	if err != nil {
		appLogger.Fatal(fmt.Errorf("run - cannot connect to db: %w", err))
	}
	defer pg.Close()

	repos, err := repo.New(pg)
	if err != nil {
		appLogger.Fatal(fmt.Errorf("run - cannot create repo: %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Services
	services, err := service.NewService(repos, appLogger)
	if err != nil {
		appLogger.Fatal(fmt.Errorf("run - cannot create services: %w", err))
	}

	go services.Polling.Run(ctx, cfg.AccrualAddress)

	// process orders that has one of status (NEW, REGISTERED, PROCESSING)
	if len(repos.NotProcessedOrders) > 0 {
		services.Polling.AddBulkToPoll("/api/orders/", repos.NotProcessedOrders)
	}

	//HTTP server
	handler := controller.NewHandler(services, appLogger)
	srv := httpserver.New(handler.InitRoutes(), httpserver.Address(cfg.Address))

	appLogger.Info("run - params: a=%s | d=%s | r=%s",
		cfg.Address, cfg.Dsn, cfg.AccrualAddress)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		appLogger.Info("run - signal: %v", s.String())
	case err := <-srv.Notify():
		appLogger.Error(fmt.Errorf("run - Notify: %w", err))
	}

	// Shutdown
	cancel()
	err = srv.Shutdown()
	if err != nil {
		appLogger.Error(fmt.Errorf("run - Shutdown: %w", err))
	}
}
