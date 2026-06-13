package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"library-api/internal/models"
	"library-api/internal/service"
)

type AuthorHandler struct {
	service *service.AuthorService
}

func NewAuthorHandler(service *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

func (h *AuthorHandler) GetAuthors(c *gin.Context) {

	authors, err := h.service.GetAuthors()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch authors",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": authors,
	})
}

func (h *AuthorHandler) GetAuthor(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid author id",
		})

		return
	}

	author, err := h.service.GetAuthor(uint(id))
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "author not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": author,
	})
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {

	var author models.Author

	if err := c.ShouldBindJSON(&author); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err := h.service.CreateAuthor(&author)
	if err != nil {

		if errors.Is(err, service.ErrInvalidAuthor) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create author",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": author,
	})
}
