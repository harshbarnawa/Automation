package gateway

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/apikey"
	"github.com/harshbarnawa/mintok/backend/internal/domain"
)

type keyStore struct{}

func (keyStore) Create(context.Context, string, string, string, string, string) (domain.APIKey, error) {
	return domain.APIKey{}, nil
}
func (keyStore) ListByProject(context.Context, string, string) ([]domain.APIKey, error) {
	return nil, nil
}
func (keyStore) Revoke(context.Context, string, string, string) error { return nil }
func (keyStore) Authenticate(context.Context, string) (string, error) { return "project-1", nil }

type engine struct{}

func (engine) Complete(_ *gin.Context, _ string, request Request) (Completion, error) {
	return Completion{ID: "chat-1", Model: request.Model, Content: "Hello", CreatedAt: time.Unix(1, 0)}, nil
}

func TestChatCompletionsReturnsOpenAICompatibleResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(apikey.NewService(keyStore{}), engine{}).Register(router)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", http.NoBody)
	req.Header.Set("Authorization", "Bearer mintok_test")
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(`{"model":"test-model","messages":[{"role":"user","content":"Hi"}]}`))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}
