
package service

import (
	"errors"

	"library-api/internal/models"
	"library-api/internal/repository"
)

var ErrInvalidAuthor = errors.New("invalid author")

type AuthorService struct {
	repo *repository.AuthorRepository
}

func NewAuthorService(repo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{
		repo: repo,
	}
}
func (s *AuthorService) CreateAuthor(author *models.Author) error {

	if author.Name == "" {
		return ErrInvalidAuthor
	}

	return s.repo.Create(author)
}

func (s *AuthorService) GetAuthor(id uint) (*models.Author, error) {
	return s.repo.GetByID(id)
}

func (s *AuthorService) GetAuthors() ([]models.Author, error) {
	return s.repo.GetAll()
}
