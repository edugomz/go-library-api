package app

import (
	"library-api/internal/middleware"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

func registerAPIroutes(r *gin.Engine, h *Handlers, authService *service.AuthService) {
	g := r.Group("/api/v1")

	// auth — public
	auth := g.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	// users — public reads, admin write (no JWT required)
	g.POST("/users", h.User.CreateUser)
	g.GET("/users", h.User.GetUsers)

	// protected write routes
	protected := g.Group("/", middleware.RequireAuth(authService))
	{
		protected.GET("/books", h.Book.GetBooks)
		protected.GET("/books/:id", h.Book.GetBook)
		protected.POST("/books", h.Book.CreateBook)
		protected.GET("/books/:id/reviews", h.Review.GetReviews)
		protected.POST("/books/:id/reviews", h.Review.CreateReview)

		protected.GET("/authors", h.Author.GetAuthors)
		protected.GET("/authors/:id", h.Author.GetAuthor)
		protected.POST("/authors", h.Author.CreateAuthor)
	}
}

func registerWebRoutes(r *gin.Engine, h *Handlers) {
	r.GET("/", h.Web.Books)
	r.GET("/books", h.Web.Books)
	r.GET("/books/:id", h.Web.BookDetail)
	r.GET("/authors", h.Web.Authors)
	r.GET("/authors/:id", h.Web.AuthorDetail)

	r.GET("/login", h.Web.LoginPage)
	r.POST("/login", h.Web.Login)
	r.GET("/register", h.Web.RegisterPage)
	r.POST("/register", h.Web.Register)
	r.POST("/logout", h.Web.Logout)
}
