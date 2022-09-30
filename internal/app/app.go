package app

import (
	"fmt"
	"github.com/PaulYakow/gophermart/config"
	ctrl "github.com/PaulYakow/gophermart/internal/controller/v1"
	"github.com/PaulYakow/gophermart/internal/pkg/httpserver"
	"github.com/PaulYakow/gophermart/internal/pkg/logger"
	"github.com/PaulYakow/gophermart/internal/usecase"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Cfg) {
	l := logger.New()
	defer l.Exit()

	// Postgres storage

	// Usecase
	someUseCase := usecase.NewServerUC()

	//HTTP server
	handler := ctrl.NewRouter(someUseCase, l)
	srv := httpserver.New(handler, httpserver.Address(cfg.Address))

	l.Info("server - run with params: a=%s | d=%s | r=%s",
		cfg.Address, cfg.Dsn, cfg.AccrualAddress)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("server - Run - signal: %v", s.String())
	case err := <-srv.Notify():
		l.Error(fmt.Errorf("server - Run - Notify: %w", err))
	}

	// Shutdown
	err := srv.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("server - Run - Shutdown: %w", err))
	}
}
