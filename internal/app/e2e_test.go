package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"library-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// e2e tests drive the real router (registerAPIroutes/middleware/handlers/
// services/repositories) over HTTP against the docker-compose.test.yml
// Postgres instance, unlike the mocked unit tests and DB-only repository
// integration tests elsewhere in the codebase.
func setupE2EServer(t *testing.T) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=postgres password=postgres dbname=library_test port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Author{},
		&models.Book{},
		&models.Review{},
		&models.ReadingList{},
	); err != nil {
		t.Fatal("migration failed:", err)
	}

	handlers := NewHandlers(db, "e2e-test-secret")

	r := gin.New()
	registerAPIroutes(r, handlers, handlers.AuthService)
	registerWebRoutes(r, handlers)

	server := httptest.NewServer(r)
	t.Cleanup(server.Close)
	return server
}

func postJSON(t *testing.T, url, token string, body map[string]any) *http.Response {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestE2E_RegisterLoginCreateBookAndReview(t *testing.T) {
	server := setupE2EServer(t)
	client := server.Client()

	email := fmt.Sprintf("e2e-%d@test.com", time.Now().UnixNano())

	// Register
	resp := postJSON(t, server.URL+"/api/v1/auth/register", "", map[string]any{
		"name":     "E2E User",
		"email":    email,
		"password": "password123",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Login
	resp = postJSON(t, server.URL+"/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": "password123",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: expected 200, got %d", resp.StatusCode)
	}
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if loginResp.Token == "" {
		t.Fatal("expected non-empty token")
	}

	// Protected routes reject requests with no token
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/books", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated books list: expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Create author (protected)
	resp = postJSON(t, server.URL+"/api/v1/authors", loginResp.Token, map[string]any{
		"name": "J.R.R. Tolkien",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create author: expected 201, got %d", resp.StatusCode)
	}
	var authorResp struct {
		Data models.Author `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authorResp); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	// Create book (protected)
	resp = postJSON(t, server.URL+"/api/v1/books", loginResp.Token, map[string]any{
		"title":     "The Hobbit",
		"author_id": authorResp.Data.ID,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create book: expected 201, got %d", resp.StatusCode)
	}
	var bookResp struct {
		Data models.Book `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&bookResp); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	// Create review (protected)
	resp = postJSON(t, fmt.Sprintf("%s/api/v1/books/%d/reviews", server.URL, bookResp.Data.ID), loginResp.Token, map[string]any{
		"rating":  5,
		"comment": "A classic",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create review: expected 201, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Fetch reviews (protected) and confirm the review round-tripped through the full stack
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/books/%d/reviews", server.URL, bookResp.Data.ID), nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get reviews: expected 200, got %d", resp.StatusCode)
	}
	var reviewsResp struct {
		Data []models.Review `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&reviewsResp); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if len(reviewsResp.Data) != 1 || reviewsResp.Data[0].Rating != 5 {
		t.Fatalf("expected 1 review with rating 5, got %+v", reviewsResp.Data)
	}
}

func TestE2E_Register_ValidationError(t *testing.T) {
	server := setupE2EServer(t)

	resp := postJSON(t, server.URL+"/api/v1/auth/register", "", map[string]any{
		"name":     "Bad User",
		"email":    "not-an-email",
		"password": "password123",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid email, got %d", resp.StatusCode)
	}
}
