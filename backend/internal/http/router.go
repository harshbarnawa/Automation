package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/auth"
	"github.com/harshbarnawa/mintok/backend/internal/config"
	applogger "github.com/harshbarnawa/mintok/backend/internal/logger"
)

// Dependencies carries the wired-up services the router mounts as routes.
// Fields may be nil, in which case their routes are not registered.
type Dependencies struct {
	Auth *auth.Handler
}

func NewRouter(cfg config.Config, log *slog.Logger, deps Dependencies) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(applogger.RequestMiddleware(log), gin.Recovery())

	if deps.Auth != nil {
		deps.Auth.Register(router)
	}

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
