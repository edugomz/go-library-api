package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"library-api/internal/service"
)

type WebHandler struct {
	bookService *service.BookService
}

func NewWebHandler(bookService *service.BookService) *WebHandler {
	return &WebHandler{
		bookService: bookService,
	}
}

func (h *WebHandler) Home(c *gin.Context) {

	books, err := h.bookService.GetBooks()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to load books")
		return
	}

	c.HTML(http.StatusOK, "books.html", gin.H{
		"books": books,
	})
}
