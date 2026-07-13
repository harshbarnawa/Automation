package gateway

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/apikey"
)

// Completion is the normalized response returned by a provider adapter.
type Completion struct {
	ID, Model, Content string
	CreatedAt          time.Time
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
type Engine interface {
	Complete(ctx *gin.Context, projectID string, request Request) (Completion, error)
}

// Handler exposes the OpenAI-compatible chat-completions transport.
type Handler struct {
	keys   *apikey.Service
	engine Engine
}

func NewHandler(keys *apikey.Service, engine Engine) *Handler {
	return &Handler{keys: keys, engine: engine}
}
func (h *Handler) Register(router gin.IRoutes) { router.POST("/v1/chat/completions", h.complete) }
func (h *Handler) complete(ctx *gin.Context) {
	key, ok := bearer(ctx.GetHeader("Authorization"))
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "invalid API key", "type": "authentication_error"}})
		return
	}
	projectID, err := h.keys.Authenticate(ctx, key)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"message": "invalid API key", "type": "authentication_error"}})
		return
	}
	var request Request
	if err := ctx.ShouldBindJSON(&request); err != nil || request.Model == "" || len(request.Messages) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "model and messages are required", "type": "invalid_request_error"}})
		return
	}
	completion, err := h.engine.Complete(ctx, projectID, request)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": gin.H{"message": "no provider is configured for this project", "type": "provider_error"}})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": completion.ID, "object": "chat.completion", "created": completion.CreatedAt.Unix(), "model": completion.Model, "choices": []gin.H{{"index": 0, "message": gin.H{"role": "assistant", "content": completion.Content}, "finish_reason": "stop"}}})
}
func bearer(header string) (string, bool) {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	value := strings.TrimSpace(header[len(prefix):])
	return value, value != ""
}
