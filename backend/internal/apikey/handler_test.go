package apikey

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harshbarnawa/mintok/backend/internal/auth"
)

func TestRoutesIssueListAndRevokeAPIKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &fakeStore{}
	service := NewService(store)
	service.rand = func(bytes []byte) (int, error) {
		for i := range bytes {
			bytes[i] = 1
		}
		return len(bytes), nil
	}
	tokens := auth.NewTokenManager("secret", 0, "test")
	handler := NewHandler(service)
	router := gin.New()
	handler.Register(router, tokens)

	token, err := tokens.Generate("user-1", "user@example.com")
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	create := httptest.NewRequest(http.MethodPost, "/projects/project-1/api-keys", bytes.NewBufferString(`{"name":"Gateway"}`))
	create.Header.Set("Authorization", "Bearer "+token)
	create.Header.Set("Content-Type", "application/json")
	created := httptest.NewRecorder()
	router.ServeHTTP(created, create)
	if created.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", created.Code, created.Body.String())
	}
	var body struct {
		Key       string `json:"key"`
		KeyPrefix string `json:"key_prefix"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Key == "" || body.KeyPrefix == "" {
		t.Fatal("expected one-time key and prefix")
	}

	list := httptest.NewRequest(http.MethodGet, "/projects/project-1/api-keys", nil)
	list.Header.Set("Authorization", "Bearer "+token)
	listed := httptest.NewRecorder()
	router.ServeHTTP(listed, list)
	if listed.Code != http.StatusOK || bytes.Contains(listed.Body.Bytes(), []byte(body.Key)) {
		t.Fatalf("expected metadata-only listing, got %d: %s", listed.Code, listed.Body.String())
	}

	revoke := httptest.NewRequest(http.MethodDelete, "/projects/project-1/api-keys/key-1", nil)
	revoke.Header.Set("Authorization", "Bearer "+token)
	revoked := httptest.NewRecorder()
	router.ServeHTTP(revoked, revoke)
	if revoked.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", revoked.Code)
	}
}

func TestRoutesRequireDashboardAccessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(NewService(&fakeStore{}))
	handler.Register(router, auth.NewTokenManager("secret", 0, "test"))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/projects/project-1/api-keys", nil))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}
