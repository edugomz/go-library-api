package app

import (
	"library-api/internal/handlers/api"
	"library-api/internal/handlers/web"
	"library-api/internal/repository"
	"library-api/internal/service"

	"gorm.io/gorm"
)

type Handlers struct {
	Auth        *api.AuthHandler
	AuthService *service.AuthService
	User        *api.UserHandler
	Book        *api.BookHandler
	Author      *api.AuthorHandler
	Review      *api.ReviewHandler

	Web *web.WebHandler
}

func NewHandlers(db *gorm.DB, jwtSecret string) *Handlers {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)

	return &Handlers{
		Auth:        api.NewAuthHandler(authService),
		AuthService: authService,
		User: api.NewUserHandler(
			service.NewUserService(userRepo),
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

		Review: api.NewReviewHandler(
			service.NewReviewService(
				repository.NewReviewRepository(db),
			),
		),

		Web: web.NewWebHandler(
			service.NewBookService(
				repository.NewBookRepository(db),
			),
		),
	}
}
