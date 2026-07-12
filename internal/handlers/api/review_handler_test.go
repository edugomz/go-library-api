package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

type mockReviewRepo struct {
	reviews []models.Review
}

func (m *mockReviewRepo) Create(review *models.Review) error {
	m.reviews = append(m.reviews, *review)
	return nil
}

func (m *mockReviewRepo) GetByBookID(bookID uint) ([]models.Review, error) {
	var out []models.Review
	for _, r := range m.reviews {
		if r.BookID == bookID {
			out = append(out, r)
		}
	}
	return out, nil
}

// withFakeUser simulates the auth middleware having already run and set userID.
func withFakeUser(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestReviewHandler_CreateReview_Success(t *testing.T) {
	repo := &mockReviewRepo{}
	h := NewReviewHandler(service.NewReviewService(repo))

	r := gin.New()
	r.POST("/books/:id/reviews", withFakeUser(1), h.CreateReview)

	req := httptest.NewRequest(http.MethodPost, "/books/5/reviews", bytes.NewBufferString(`{"rating":4,"comment":"good"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	if len(repo.reviews) != 1 || repo.reviews[0].BookID != 5 || repo.reviews[0].UserID != 1 {
		t.Fatalf("unexpected stored review: %+v", repo.reviews)
	}
}

func TestReviewHandler_CreateReview_InvalidBookID(t *testing.T) {
	h := NewReviewHandler(service.NewReviewService(&mockReviewRepo{}))

	r := gin.New()
	r.POST("/books/:id/reviews", withFakeUser(1), h.CreateReview)

	req := httptest.NewRequest(http.MethodPost, "/books/abc/reviews", bytes.NewBufferString(`{"rating":4}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestReviewHandler_CreateReview_InvalidRating(t *testing.T) {
	h := NewReviewHandler(service.NewReviewService(&mockReviewRepo{}))

	r := gin.New()
	r.POST("/books/:id/reviews", withFakeUser(1), h.CreateReview)

	req := httptest.NewRequest(http.MethodPost, "/books/5/reviews", bytes.NewBufferString(`{"rating":9}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestReviewHandler_GetReviews(t *testing.T) {
	repo := &mockReviewRepo{reviews: []models.Review{{BookID: 5, Rating: 4}}}
	h := NewReviewHandler(service.NewReviewService(repo))

	r := gin.New()
	r.GET("/books/:id/reviews", h.GetReviews)

	req := httptest.NewRequest(http.MethodGet, "/books/5/reviews", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestReviewHandler_GetReviews_InvalidBookID(t *testing.T) {
	h := NewReviewHandler(service.NewReviewService(&mockReviewRepo{}))

	r := gin.New()
	r.GET("/books/:id/reviews", h.GetReviews)

	req := httptest.NewRequest(http.MethodGet, "/books/abc/reviews", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
