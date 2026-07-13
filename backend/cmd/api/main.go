package main

import (
	"log"

	"github.com/harshbarnawa/mintok/backend/internal/config"
	"github.com/harshbarnawa/mintok/backend/internal/http"
)

func main() {
	cfg := config.Load()
	router := http.NewRouter(cfg)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start api: %v", err)
	}
}
