package service

import (
	"errors"
	"library-api/internal/models"
	"library-api/internal/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}
func (s *BookService) GetBooks() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetBook(id uint) (*models.Book, error) {
	return s.repo.GetByID(id)
}


var ErrInvalidBook = errors.New("invalid book data")

func (s *BookService) CreateBook(book *models.Book) error {

	// business layer is where validation would go later
	// e.g. check author exists, title not empty, etc.

	if book.Title == "" {
		return ErrInvalidBook
	}

	return s.repo.Create(book)
}

