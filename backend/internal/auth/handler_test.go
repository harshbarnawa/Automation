package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestHandler() *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := NewService(newFakeUserRepo(), NewBcryptHasher(4))
	handler := NewHandler(svc, NewTokenManager("test-secret", 0, "mintok-test"), NewRefreshTokenManager(newFakeRefreshTokenStore(), 0))

	router := gin.New()
	handler.Register(router)
	return router
}

func doJSON(t *testing.T, router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestRegisterEndpoint(t *testing.T) {
	router := newTestHandler()

	recorder := doJSON(t, router, http.MethodPost, "/auth/register", gin.H{
		"email":    "new@example.com",
		"name":     "New User",
		"password": "supersecret",
	})

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var body struct {
		User         UserResponse `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.User.Email != "new@example.com" {
		t.Fatalf("unexpected email %q", body.User.Email)
	}
	if body.User.ID == "" {
		t.Fatal("expected user id in response")
	}
	if body.AccessToken == "" {
		t.Fatal("expected access token in response")
	}
	if body.RefreshToken == "" {
		t.Fatal("expected refresh token in response")
	}
}

func TestRefreshEndpointRotatesToken(t *testing.T) {
	router := newTestHandler()
	registration := doJSON(t, router, http.MethodPost, "/auth/register", gin.H{
		"email": "refresh@example.com", "name": "Refresh", "password": "supersecret",
	})

	var registered struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(registration.Body.Bytes(), &registered); err != nil {
		t.Fatalf("decode registration: %v", err)
	}

	refresh := doJSON(t, router, http.MethodPost, "/auth/refresh", gin.H{"refresh_token": registered.RefreshToken})
	if refresh.Code != http.StatusOK {
		t.Fatalf("expected refresh 200, got %d: %s", refresh.Code, refresh.Body.String())
	}

	if reused := doJSON(t, router, http.MethodPost, "/auth/refresh", gin.H{"refresh_token": registered.RefreshToken}); reused.Code != http.StatusUnauthorized {
		t.Fatalf("expected reused token 401, got %d", reused.Code)
	}
}

func TestRegisterEndpointDuplicate(t *testing.T) {
	router := newTestHandler()
	payload := gin.H{"email": "dup@example.com", "name": "Dup", "password": "supersecret"}

	if code := doJSON(t, router, http.MethodPost, "/auth/register", payload).Code; code != http.StatusCreated {
		t.Fatalf("expected first register 201, got %d", code)
	}

	recorder := doJSON(t, router, http.MethodPost, "/auth/register", payload)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", recorder.Code)
	}
}

func TestRegisterEndpointValidation(t *testing.T) {
	router := newTestHandler()

	recorder := doJSON(t, router, http.MethodPost, "/auth/register", gin.H{
		"email":    "a@b.com",
		"name":     "A",
		"password": "short",
	})
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestLoginEndpoint(t *testing.T) {
	router := newTestHandler()
	register := gin.H{"email": "login@example.com", "name": "Login", "password": "supersecret"}
	if code := doJSON(t, router, http.MethodPost, "/auth/register", register).Code; code != http.StatusCreated {
		t.Fatalf("expected register 201, got %d", code)
	}

	recorder := doJSON(t, router, http.MethodPost, "/auth/login", gin.H{
		"email":    "login@example.com",
		"password": "supersecret",
	})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AccessToken == "" {
		t.Fatal("expected access token in response")
	}
}

func TestMeEndpoint(t *testing.T) {
	router := newTestHandler()
	register := gin.H{"email": "me@example.com", "name": "Me", "password": "supersecret"}
	registration := doJSON(t, router, http.MethodPost, "/auth/register", register)

	var body struct {
		User        UserResponse `json:"user"`
		AccessToken string       `json:"access_token"`
	}
	if err := json.Unmarshal(registration.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode registration: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	request.Header.Set("Authorization", "Bearer "+body.AccessToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestMeEndpointRequiresToken(t *testing.T) {
	router := newTestHandler()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/auth/me", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestLoginEndpointInvalid(t *testing.T) {
	router := newTestHandler()

	recorder := doJSON(t, router, http.MethodPost, "/auth/login", gin.H{
		"email":    "missing@example.com",
		"password": "supersecret",
	})
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}
