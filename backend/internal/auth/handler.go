package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/domain"
)

// Handler exposes the authentication HTTP endpoints.
type Handler struct {
	service       *Service
	tokens        *TokenManager
	refreshTokens *RefreshTokenManager
}

// NewHandler builds an auth Handler.
func NewHandler(service *Service, tokens *TokenManager, refreshTokens *RefreshTokenManager) *Handler {
	return &Handler{service: service, tokens: tokens, refreshTokens: refreshTokens}
}

// Register wires the auth routes onto the given router group.
func (h *Handler) Register(group gin.IRoutes) {
	group.POST("/auth/register", h.register)
	group.POST("/auth/login", h.login)
	group.POST("/auth/refresh", h.refresh)
	group.GET("/auth/me", Middleware(h.tokens), h.me)
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

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// UserResponse is the public JSON representation of a user.
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type authenticationResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
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

	h.writeAuthenticationResponse(ctx, http.StatusCreated, user)
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

	h.writeAuthenticationResponse(ctx, http.StatusOK, user)
}

func (h *Handler) me(ctx *gin.Context) {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (h *Handler) refresh(ctx *gin.Context) {
	var req refreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, refreshToken, err := h.refreshTokens.Rotate(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	accessToken, err := h.tokens.Generate(userID, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    int64(h.tokens.AccessTTL().Seconds()),
	})
}

func (h *Handler) writeAuthenticationResponse(ctx *gin.Context, status int, user domain.User) {
	accessToken, err := h.tokens.Generate(user.ID, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	refreshToken, err := h.refreshTokens.Issue(ctx.Request.Context(), user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(status, authenticationResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(h.tokens.AccessTTL().Seconds()),
	})
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
