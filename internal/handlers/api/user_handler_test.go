package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

type mockUserRepo struct {
	users     []models.User
	createErr error
}

func (m *mockUserRepo) Create(user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users = append(m.users, *user)
	return nil
}

func (m *mockUserRepo) GetAll() ([]models.User, error) {
	return m.users, nil
}

func init() {
	gin.SetMode(gin.TestMode)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	repo := &mockUserRepo{}
	h := NewUserHandler(service.NewUserService(repo))

	r := gin.New()
	r.POST("/users", h.CreateUser)

	body := `{"name":"John","email":"john@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	repo := &mockUserRepo{}
	h := NewUserHandler(service.NewUserService(repo))

	r := gin.New()
	r.POST("/users", h.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	repo := &mockUserRepo{createErr: errors.New("db error")}
	h := NewUserHandler(service.NewUserService(repo))

	r := gin.New()
	r.POST("/users", h.CreateUser)

	body := `{"name":"John","email":"john@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	repo := &mockUserRepo{users: []models.User{{ID: 1, Name: "John", Email: "john@test.com"}}}
	h := NewUserHandler(service.NewUserService(repo))

	r := gin.New()
	r.GET("/users", h.GetUsers)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Data []models.User `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 user, got %d", len(resp.Data))
	}
}
