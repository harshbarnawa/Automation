package main

import (
	"context"
	"time"

	"github.com/harshbarnawa/mintok/backend/internal/config"
	"github.com/harshbarnawa/mintok/backend/internal/database"
	"github.com/harshbarnawa/mintok/backend/internal/http"
	"github.com/harshbarnawa/mintok/backend/internal/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg)
	router := http.NewRouter(cfg, log)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewPool(ctx, cfg)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		panic(err)
	}
	defer db.Close()

	if err := database.ApplyMigrations(ctx, db, "migrations"); err != nil {
		log.Error("failed to apply database migrations", "error", err)
		panic(err)
	}

	log.Info("starting api", "port", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("failed to start api", "error", err)
		panic(err)
	}
}
