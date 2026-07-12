package service

import (
	"errors"
	"testing"
	"time"

	"library-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthUserRepository struct {
	users     map[string]models.User
	createErr error
}

func newMockAuthUserRepository() *MockAuthUserRepository {
	return &MockAuthUserRepository{users: map[string]models.User{}}
}

func (m *MockAuthUserRepository) Create(user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.Email] = *user
	return nil
}

func (m *MockAuthUserRepository) FindByEmail(email string) (*models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return &user, nil
}

func TestAuthService_Register(t *testing.T) {
	repo := newMockAuthUserRepository()
	svc := NewAuthService(repo, "test-secret")

	if err := svc.Register("Jane", "jane@test.com", "password123"); err != nil {
		t.Fatal(err)
	}

	stored, ok := repo.users["jane@test.com"]
	if !ok {
		t.Fatal("expected user to be stored")
	}

	if stored.Password == "password123" {
		t.Fatal("expected password to be hashed, not stored in plaintext")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte("password123")); err != nil {
		t.Fatalf("stored hash does not match password: %v", err)
	}
}

func TestAuthService_Register_RepositoryError(t *testing.T) {
	repo := newMockAuthUserRepository()
	repo.createErr = errors.New("db error")
	svc := NewAuthService(repo, "test-secret")

	if err := svc.Register("Jane", "jane@test.com", "password123"); err == nil {
		t.Fatal("expected error from repository to propagate")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := newMockAuthUserRepository()
	svc := NewAuthService(repo, "test-secret")

	if err := svc.Register("Jane", "jane@test.com", "password123"); err != nil {
		t.Fatal(err)
	}

	token, err := svc.Login("jane@test.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	userID, err := svc.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}

	stored := repo.users["jane@test.com"]
	if userID != stored.ID {
		t.Fatalf("expected userID %d, got %d", stored.ID, userID)
	}
}

func TestAuthService_Login_UnknownEmail(t *testing.T) {
	repo := newMockAuthUserRepository()
	svc := NewAuthService(repo, "test-secret")

	if _, err := svc.Login("missing@test.com", "password123"); err == nil {
		t.Fatal("expected error for unknown email")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	repo := newMockAuthUserRepository()
	svc := NewAuthService(repo, "test-secret")

	if err := svc.Register("Jane", "jane@test.com", "password123"); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.Login("jane@test.com", "wrongpassword"); err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestAuthService_ParseToken_Invalid(t *testing.T) {
	svc := NewAuthService(newMockAuthUserRepository(), "test-secret")

	if _, err := svc.ParseToken("not-a-real-token"); err == nil {
		t.Fatal("expected error for malformed token")
	}
}

func TestAuthService_ParseToken_WrongSecret(t *testing.T) {
	repo := newMockAuthUserRepository()
	svc := NewAuthService(repo, "test-secret")

	if err := svc.Register("Jane", "jane@test.com", "password123"); err != nil {
		t.Fatal(err)
	}

	token, err := svc.Login("jane@test.com", "password123")
	if err != nil {
		t.Fatal(err)
	}

	other := NewAuthService(repo, "different-secret")
	if _, err := other.ParseToken(token); err == nil {
		t.Fatal("expected error when validating with wrong secret")
	}
}

func TestAuthService_ParseToken_Expired(t *testing.T) {
	secret := "test-secret"
	claims := jwt.MapClaims{
		"sub": float64(1),
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	svc := NewAuthService(newMockAuthUserRepository(), secret)
	if _, err := svc.ParseToken(signed); err == nil {
		t.Fatal("expected error for expired token")
	}
}
