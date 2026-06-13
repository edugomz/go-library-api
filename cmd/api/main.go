// package main
//
// import (
// 	"library-api/internal/db"
// 	"library-api/internal/handlers"
// 	"library-api/internal/models"
// 	"library-api/internal/service"
//
// 	"github.com/gin-gonic/gin"
// )
//
// func main() {
// 	dsn := "host=localhost user=postgres password=postgres dbname=library port=5432 sslmode=disable"
// 	db.Connect(dsn)
// 	db.DB.AutoMigrate(
// 		&models.User{},
// 		&models.Author{},
// 		&models.Book{},
// 		&models.Review{},
// 		&models.ReadingList{},
// 	)
//
// 	r := gin.Default()
//
// 	// userHandler := handlers.NewUserHandler(userService)
// 	userHandler := handlers.UserHandler{
// 		service: &service.UserService
// 	}
//
// 	r.POST("/users", userHandler.CreateUser)
// 	r.GET("/users", userHandler.GetUsers)
//
// 	r.GET("/health", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"status": "ok"})
// 	})
//
// 	r.Run(":8080")
// }

package main

import (
	"log"
	// "log/slog"
	// "os"


	"github.com/gin-gonic/gin"

	"library-api/internal/db"
	"library-api/internal/handlers/web"
	"library-api/internal/handlers/api"
	"library-api/internal/logger"
	"library-api/internal/models"
	"library-api/internal/repository"
	"library-api/internal/service"
	"library-api/internal/config"
)

func main() {


	// =========================
	// 1. Initialize DB
	// =========================
	// dsn := "host=localhost user=postgres password=postgres dbname=library port=5432 sslmode=disable"
  //
	// if err := db.Connect(dsn); err != nil {
	// 	logger.Log.Fatal("failed to connect database:", err)
	// }

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Connect(cfg.DSN()); err != nil {
		logger.Log.Error("error", err.Error())
	}

	// Auto-migrate (dev only)
	if err := db.DB.AutoMigrate(
		&models.User{},
		&models.Author{},
		&models.Book{},
		&models.Review{},
		&models.ReadingList{},
	); err != nil {
		log.Fatal("migration failed:", err)
	}

	// =========================
	// 2. Initialize layers
	// =========================

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)

	bookRepo := repository.NewBookRepository(db.DB)
	bookService := service.NewBookService(bookRepo)
	bookHandler := api.NewBookHandler(bookService)

	authorRepo := repository.NewAuthorRepository(db.DB)
	authorService := service.NewAuthorService(authorRepo)
	authorHandler := api.NewAuthorHandler(authorService)

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// =========================
	// 4. Register routes
	// =========================

	apiV1 := r.Group("/api/v1")

	{
		apiV1.POST("/users", userHandler.CreateUser)
		apiV1.GET("/users", userHandler.GetUsers)

		apiV1.GET("/books", bookHandler.GetBooks)
		apiV1.GET("/books/:id", bookHandler.GetBook)
		apiV1.POST("/books", bookHandler.CreateBook)

		apiV1.GET("/authors", authorHandler.GetAuthors)
		apiV1.GET("/authors/:id", authorHandler.GetAuthor)
		apiV1.POST("/authors", authorHandler.CreateAuthor)
	}

	r.LoadHTMLGlob("internal/views/*")
	webHandler := web.NewWebHandler(bookService)

	r.GET("/", webHandler.Home)


	// =========================
	// 5. Start server
	// =========================
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
