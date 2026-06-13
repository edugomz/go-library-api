package service

import (
    "testing"

    "library-api/internal/models"
)

type MockUserRepository struct {
    users []models.User
}

func (m *MockUserRepository) Create(user *models.User) error {
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
                ID: 1,
                Name: "John",
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
