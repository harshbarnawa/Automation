package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/config"
)

func NewRouter(cfg config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service":     cfg.ServiceName,
			"environment": cfg.Environment,
			"status":      "ok",
		})
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"name":   "Mintok",
			"status": "under_active_development",
		})
	})

	return router
}
