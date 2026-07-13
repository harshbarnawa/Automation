package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/domain"
)

// Handler exposes the authentication HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler builds an auth Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register wires the auth routes onto the given router group.
func (h *Handler) Register(group gin.IRoutes) {
	group.POST("/auth/register", h.register)
	group.POST("/auth/login", h.login)
}

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse is the public JSON representation of a user.
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NewUserResponse converts a domain user into its public representation.
func NewUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *Handler) register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.Register(ctx.Request.Context(), RegisterParams{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		writeAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": NewUserResponse(user)})
}

func (h *Handler) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.Login(ctx.Request.Context(), LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeAuthError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": NewUserResponse(user)})
}

func writeAuthError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrEmailTaken):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidCredentials):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
