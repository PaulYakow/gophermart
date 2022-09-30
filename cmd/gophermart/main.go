package main

import (
	"github.com/PaulYakow/gophermart/config"
	"github.com/PaulYakow/gophermart/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("create config: %v\n", err)
	}

	app.Run(cfg)
}
