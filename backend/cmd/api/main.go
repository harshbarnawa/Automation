package main

import (
	"context"
	"time"

	"github.com/harshbarnawa/mintok/backend/internal/auth"
	"github.com/harshbarnawa/mintok/backend/internal/cache"
	"github.com/harshbarnawa/mintok/backend/internal/config"
	"github.com/harshbarnawa/mintok/backend/internal/database"
	"github.com/harshbarnawa/mintok/backend/internal/http"
	"github.com/harshbarnawa/mintok/backend/internal/logger"
	"github.com/harshbarnawa/mintok/backend/internal/repository"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewPool(ctx, cfg)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		panic(err)
	}
	defer db.Close()

	redisClient, err := cache.NewRedisClient(ctx, cfg)
	if err != nil {
		log.Error("failed to connect redis", "error", err)
		panic(err)
	}
	defer func() {
		_ = redisClient.Close()
	}()

	if err := database.ApplyMigrations(ctx, db, "migrations"); err != nil {
		log.Error("failed to apply database migrations", "error", err)
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := auth.NewService(userRepo, auth.NewBcryptHasher(0))
	authHandler := auth.NewHandler(authService)

	router := http.NewRouter(cfg, log, http.Dependencies{
		Auth: authHandler,
	})

	log.Info("starting api", "port", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("failed to start api", "error", err)
		panic(err)
	}
}
