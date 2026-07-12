package repository

import (
	"fmt"
	"testing"

	"library-api/internal/models"
)

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Name:  "John",
		Email: fmt.Sprintf("john-%d@test.com", uniqueSuffix()),
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Fatal("expected created user to have an ID")
	}

	users, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Fatal("expected users")
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	email := fmt.Sprintf("jane-%d@test.com", uniqueSuffix())
	user := &models.User{Name: "Jane", Email: email}
	if err := repo.Create(user); err != nil {
		t.Fatal(err)
	}

	found, err := repo.FindByEmail(email)
	if err != nil {
		t.Fatal(err)
	}

	if found.ID != user.ID {
		t.Fatalf("expected user ID %d, got %d", user.ID, found.ID)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByEmail(fmt.Sprintf("missing-%d@test.com", uniqueSuffix()))
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}
