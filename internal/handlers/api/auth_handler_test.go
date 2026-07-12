package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

type mockAuthUserRepo struct {
	users map[string]models.User
}

func newMockAuthUserRepo() *mockAuthUserRepo {
	return &mockAuthUserRepo{users: map[string]models.User{}}
}

func (m *mockAuthUserRepo) Create(user *models.User) error {
	m.users[user.Email] = *user
	return nil
}

func (m *mockAuthUserRepo) FindByEmail(email string) (*models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return &user, nil
}

func TestAuthHandler_Register_Success(t *testing.T) {
	h := NewAuthHandler(service.NewAuthService(newMockAuthUserRepo(), "test-secret"))

	r := gin.New()
	r.POST("/auth/register", h.Register)

	body := `{"name":"Jane","email":"jane@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Register_InvalidEmail(t *testing.T) {
	h := NewAuthHandler(service.NewAuthService(newMockAuthUserRepo(), "test-secret"))

	r := gin.New()
	r.POST("/auth/register", h.Register)

	body := `{"name":"Jane","email":"not-an-email","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	h := NewAuthHandler(service.NewAuthService(newMockAuthUserRepo(), "test-secret"))

	r := gin.New()
	r.POST("/auth/register", h.Register)

	body := `{"name":"Jane","email":"jane@test.com","password":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	repo := newMockAuthUserRepo()
	authService := service.NewAuthService(repo, "test-secret")
	if err := authService.Register("Jane", "jane@test.com", "password123"); err != nil {
		t.Fatal(err)
	}

	h := NewAuthHandler(authService)

	r := gin.New()
	r.POST("/auth/login", h.Login)

	body := `{"email":"jane@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	h := NewAuthHandler(service.NewAuthService(newMockAuthUserRepo(), "test-secret"))

	r := gin.New()
	r.POST("/auth/login", h.Login)

	body := `{"email":"jane@test.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
