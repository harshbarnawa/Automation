package main

import (
	"github.com/harshbarnawa/mintok/backend/internal/config"
	"github.com/harshbarnawa/mintok/backend/internal/http"
	"github.com/harshbarnawa/mintok/backend/internal/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg)
	router := http.NewRouter(cfg, log)

	log.Info("starting api", "port", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("failed to start api", "error", err)
		panic(err)
	}
}
