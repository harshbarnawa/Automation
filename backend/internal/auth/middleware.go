package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context keys under which the middleware stores the authenticated identity.
const (
	contextUserID = "auth.user_id"
	contextEmail  = "auth.email"
)

// Middleware validates the bearer access token and stores the identity on the
// request context, aborting with 401 when the token is missing or invalid.
func Middleware(tokens *TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		tokenString, ok := bearerToken(header)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := tokens.Verify(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		ctx.Set(contextUserID, claims.UserID)
		ctx.Set(contextEmail, claims.Email)
		ctx.Next()
	}
}

// UserIDFromContext returns the authenticated user id set by Middleware.
func UserIDFromContext(ctx *gin.Context) (string, bool) {
	value, ok := ctx.Get(contextUserID)
	if !ok {
		return "", false
	}
	id, ok := value.(string)
	return id, ok && id != ""
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	token := strings.TrimSpace(header[len(prefix):])
	return token, token != ""
}
