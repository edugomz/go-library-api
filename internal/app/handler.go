package app

import (
	"library-api/internal/handlers/api"
	"library-api/internal/handlers/web"
	"library-api/internal/repository"
	"library-api/internal/service"

	"gorm.io/gorm"
)

type Handlers struct {
	User *api.UserHandler
	Book *api.BookHandler
	Author *api.AuthorHandler

	Web *web.WebHandler
}

func NewHandlers(db *gorm.DB) *Handlers {
	return &Handlers{
		User: api.NewUserHandler(
			service.NewUserService(
				repository.NewUserRepository(db),
			),
		),
		Book: api.NewBookHandler(
			service.NewBookService(
				repository.NewBookRepository(db),
			),
		),
		Author: api.NewAuthorHandler(
			service.NewAuthorService(
				repository.NewAuthorRepository(db),
			),
		),

		Web: web.NewWebHandler(
			service.NewBookService(
				repository.NewBookRepository(db),
			),
		),
	}
}
