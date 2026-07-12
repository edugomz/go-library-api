package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"library-api/internal/models"
	"library-api/internal/service"
)

const tokenCookie = "token"

type WebHandler struct {
	bookService   *service.BookService
	authorService *service.AuthorService
	authService   *service.AuthService
}

func NewWebHandler(bookService *service.BookService, authorService *service.AuthorService, authService *service.AuthService) *WebHandler {
	return &WebHandler{
		bookService:   bookService,
		authorService: authorService,
		authService:   authService,
	}
}

func (h *WebHandler) isLoggedIn(c *gin.Context) bool {
	_, err := c.Cookie(tokenCookie)
	return err == nil
}

func (h *WebHandler) Books(c *gin.Context) {
	books, err := h.bookService.GetBooks()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to load books")
		return
	}

	c.HTML(http.StatusOK, "books.html", gin.H{
		"books":    books,
		"LoggedIn": h.isLoggedIn(c),
	})
}

func (h *WebHandler) BookDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid book id")
		return
	}

	book, err := h.bookService.GetBook(uint(id))
	if err != nil {
		c.String(http.StatusNotFound, "book not found")
		return
	}

	author, _ := h.authorService.GetAuthor(book.AuthorID)

	c.HTML(http.StatusOK, "book_detail.html", gin.H{
		"book":     book,
		"author":   author,
		"LoggedIn": h.isLoggedIn(c),
	})
}

func (h *WebHandler) Authors(c *gin.Context) {
	authors, err := h.authorService.GetAuthors()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to load authors")
		return
	}

	c.HTML(http.StatusOK, "authors.html", gin.H{
		"authors":  authors,
		"LoggedIn": h.isLoggedIn(c),
	})
}

func (h *WebHandler) AuthorDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid author id")
		return
	}

	author, err := h.authorService.GetAuthor(uint(id))
	if err != nil {
		c.String(http.StatusNotFound, "author not found")
		return
	}

	allBooks, err := h.bookService.GetBooks()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to load books")
		return
	}

	var books []models.Book
	for _, b := range allBooks {
		if b.AuthorID == author.ID {
			books = append(books, b)
		}
	}

	c.HTML(http.StatusOK, "author_detail.html", gin.H{
		"author":   author,
		"books":    books,
		"LoggedIn": h.isLoggedIn(c),
	})
}

func (h *WebHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"LoggedIn": h.isLoggedIn(c)})
}

func (h *WebHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	token, err := h.authService.Login(email, password)
	if err != nil {
		h.loginError(c, http.StatusUnauthorized, "Invalid email or password.")
		return
	}

	c.SetCookie(tokenCookie, token, 3600*24, "/", "", false, true)
	h.redirect(c, "/books")
}

func (h *WebHandler) loginError(c *gin.Context, status int, msg string) {
	if c.GetHeader("HX-Request") == "true" {
		c.String(status, msg)
		return
	}
	c.HTML(status, "login.html", gin.H{"error": msg, "LoggedIn": h.isLoggedIn(c)})
}

func (h *WebHandler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"LoggedIn": h.isLoggedIn(c)})
}

func (h *WebHandler) Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if err := h.authService.Register(name, email, password); err != nil {
		h.registerError(c, http.StatusConflict, "Could not register: "+err.Error())
		return
	}

	h.redirect(c, "/login")
}

func (h *WebHandler) registerError(c *gin.Context, status int, msg string) {
	if c.GetHeader("HX-Request") == "true" {
		c.String(status, msg)
		return
	}
	c.HTML(status, "register.html", gin.H{"error": msg, "LoggedIn": h.isLoggedIn(c)})
}

func (h *WebHandler) Logout(c *gin.Context) {
	c.SetCookie(tokenCookie, "", -1, "/", "", false, true)
	h.redirect(c, "/books")
}

// redirect sends an HX-Redirect for htmx requests (which don't follow a plain
// 302 as a page navigation) and a normal redirect otherwise.
func (h *WebHandler) redirect(c *gin.Context, path string) {
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", path)
		c.Status(http.StatusOK)
		return
	}
	c.Redirect(http.StatusFound, path)
}
