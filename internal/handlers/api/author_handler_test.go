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

type mockAuthorRepo struct {
	authors   []models.Author
	createErr error
}

func (m *mockAuthorRepo) GetAll() ([]models.Author, error) {
	return m.authors, nil
}

func (m *mockAuthorRepo) GetByID(id uint) (*models.Author, error) {
	for _, a := range m.authors {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockAuthorRepo) Create(author *models.Author) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.authors = append(m.authors, *author)
	return nil
}

func TestAuthorHandler_CreateAuthor_Success(t *testing.T) {
	h := NewAuthorHandler(service.NewAuthorService(&mockAuthorRepo{}))

	r := gin.New()
	r.POST("/authors", h.CreateAuthor)

	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewBufferString(`{"name":"Tolkien"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthorHandler_CreateAuthor_InvalidName(t *testing.T) {
	h := NewAuthorHandler(service.NewAuthorService(&mockAuthorRepo{}))

	r := gin.New()
	r.POST("/authors", h.CreateAuthor)

	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewBufferString(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthorHandler_GetAuthor_NotFound(t *testing.T) {
	h := NewAuthorHandler(service.NewAuthorService(&mockAuthorRepo{}))

	r := gin.New()
	r.GET("/authors/:id", h.GetAuthor)

	req := httptest.NewRequest(http.MethodGet, "/authors/42", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestAuthorHandler_GetAuthor_InvalidID(t *testing.T) {
	h := NewAuthorHandler(service.NewAuthorService(&mockAuthorRepo{}))

	r := gin.New()
	r.GET("/authors/:id", h.GetAuthor)

	req := httptest.NewRequest(http.MethodGet, "/authors/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthorHandler_GetAuthors(t *testing.T) {
	repo := &mockAuthorRepo{authors: []models.Author{{ID: 1, Name: "Tolkien"}}}
	h := NewAuthorHandler(service.NewAuthorService(repo))

	r := gin.New()
	r.GET("/authors", h.GetAuthors)

	req := httptest.NewRequest(http.MethodGet, "/authors", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
