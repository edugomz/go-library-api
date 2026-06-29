
package service

import (
	"errors"

	"library-api/internal/models"
)

var ErrInvalidAuthor = errors.New("invalid author")

type AuthorRepository interface {
	GetAll() ([]models.Author, error)
	GetByID(id uint) (*models.Author, error)
	Create(author *models.Author) error
}

type AuthorService struct {
	repo AuthorRepository
}

func NewAuthorService(repo AuthorRepository) *AuthorService {
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
