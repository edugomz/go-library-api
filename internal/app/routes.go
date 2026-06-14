package app

import "github.com/gin-gonic/gin"

func registerAPIroutes(r *gin.Engine, h *Handlers) {
	g := r.Group("/api/v1")
	{
		g.POST("/users", h.User.CreateUser)
		g.GET("/users", h.User.GetUsers)

		g.GET("/books", h.Book.GetBooks)
		g.GET("/books/:id", h.Book.GetBook)
		g.POST("/books", h.Book.CreateBook)

		g.GET("/authors", h.Author.GetAuthors)
		g.GET("/authors/:id", h.Author.GetAuthor)
		g.POST("/authors", h.Author.CreateAuthor)
	}

}

func registerWebRoutes(r *gin.Engine, h *Handlers) {
	r.GET("/", h.Web.Home)
}
