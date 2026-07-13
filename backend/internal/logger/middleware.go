package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		log.Info(
			"http_request",
			"method", ctx.Request.Method,
			"path", ctx.FullPath(),
			"status", ctx.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", ctx.ClientIP(),
		)
	}
}
