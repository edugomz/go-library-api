package service

import (
	"library-api/internal/models"
)

// type UserRepository = repository.UserRepository

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserRepository interface {
    Create(user *models.User) error
    GetAll() ([]models.User, error)
}

type UserService struct {
    repo UserRepository
}

func (s *UserService) CreateUser(name, email string) error {
	user := &models.User{Name: name, Email: email}
	return s.repo.Create(user)
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repo.GetAll()
}
