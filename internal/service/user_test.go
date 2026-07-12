package service

import (
	"errors"
	"testing"

	"library-api/internal/models"
)

type MockUserRepository struct {
	users     []models.User
	createErr error
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users = append(m.users, *user)
	return nil
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	return m.users, nil
}
func TestGetUsers(t *testing.T) {

	mockRepo := &MockUserRepository{
		users: []models.User{
			{
				ID:    1,
				Name:  "John",
				Email: "john@test.com",
			},
		},
	}

	service := NewUserService(mockRepo)

	users, err := service.GetUsers()

	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestCreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	if err := service.CreateUser("Jane", "jane@test.com"); err != nil {
		t.Fatal(err)
	}

	if len(mockRepo.users) != 1 {
		t.Fatalf("expected 1 user stored, got %d", len(mockRepo.users))
	}

	if mockRepo.users[0].Email != "jane@test.com" {
		t.Fatalf("expected email jane@test.com, got %q", mockRepo.users[0].Email)
	}
}

func TestCreateUser_RepositoryError(t *testing.T) {
	mockRepo := &MockUserRepository{createErr: errors.New("db error")}
	service := NewUserService(mockRepo)

	if err := service.CreateUser("Jane", "jane@test.com"); err == nil {
		t.Fatal("expected error from repository to propagate")
	}
}
