package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"library-api/internal/models"
	"library-api/internal/service"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}
func (h *BookHandler) GetBooks(c *gin.Context) {

	books, err := h.service.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch books",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": books,
	})
}

func (h *BookHandler) GetBook(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid book id",
		})
		return
	}

	book, err := h.service.GetBook(uint(id))
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "book not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}

func (h *BookHandler) CreateBook(c *gin.Context) {

	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err := h.service.CreateBook(&book)
	if err != nil {

		if errors.Is(err, service.ErrInvalidBook) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create book",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": book,
	})
}
