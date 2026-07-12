package middleware

import (
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

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter(authService *service.AuthService) *gin.Engine {
	r := gin.New()
	r.GET("/protected", RequireAuth(authService), func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})
	return r
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	authService := service.NewAuthService(&mockAuthUserRepo{users: map[string]models.User{}}, "test-secret")
	r := newTestRouter(authService)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_MalformedHeader(t *testing.T) {
	authService := service.NewAuthService(&mockAuthUserRepo{users: map[string]models.User{}}, "test-secret")
	r := newTestRouter(authService)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	authService := service.NewAuthService(&mockAuthUserRepo{users: map[string]models.User{}}, "test-secret")
	r := newTestRouter(authService)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_ValidToken(t *testing.T) {
	repo := &mockAuthUserRepo{users: map[string]models.User{}}
	authService := service.NewAuthService(repo, "test-secret")

	if err := authService.Register("Jane", "jane@test.com", "password123"); err != nil {
		t.Fatal(err)
	}

	token, err := authService.Login("jane@test.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	r := newTestRouter(authService)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
