package apikey

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/auth"
	"github.com/harshbarnawa/mintok/backend/internal/domain"
	"github.com/harshbarnawa/mintok/backend/internal/repository"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(router gin.IRoutes, tokens *auth.TokenManager) {
	protected := auth.Middleware(tokens)
	router.POST("/projects/:project_id/api-keys", protected, h.create)
	router.GET("/projects/:project_id/api-keys", protected, h.list)
	router.DELETE("/projects/:project_id/api-keys/:key_id", protected, h.revoke)
}

type createRequest struct {
	Name string `json:"name"`
}

type response struct {
	ID         string     `json:"id"`
	ProjectID  string     `json:"project_id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	Key        string     `json:"key,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

func (h *Handler) create(ctx *gin.Context) {
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	key, plaintext, err := h.service.Create(ctx, userID(ctx), ctx.Param("project_id"), req.Name)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, fromDomain(key, plaintext))
}

func (h *Handler) list(ctx *gin.Context) {
	keys, err := h.service.List(ctx, userID(ctx), ctx.Param("project_id"))
	if err != nil {
		writeError(ctx, err)
		return
	}
	responses := make([]response, 0, len(keys))
	for _, key := range keys {
		responses = append(responses, fromDomain(key, ""))
	}
	ctx.JSON(http.StatusOK, gin.H{"api_keys": responses})
}

func (h *Handler) revoke(ctx *gin.Context) {
	err := h.service.Revoke(ctx, userID(ctx), ctx.Param("project_id"), ctx.Param("key_id"))
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func userID(ctx *gin.Context) string { id, _ := auth.UserIDFromContext(ctx); return id }

func fromDomain(key domain.APIKey, plaintext string) response {
	return response{ID: key.ID, ProjectID: key.ProjectID, Name: key.Name, KeyPrefix: key.KeyPrefix, Key: plaintext, LastUsedAt: key.LastUsedAt, CreatedAt: key.CreatedAt, RevokedAt: key.RevokedAt}
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, repository.ErrNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "project or api key not found"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
