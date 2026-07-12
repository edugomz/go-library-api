package app

import (
	"library-api/internal/config"
	"library-api/internal/db"
	"library-api/internal/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	DB       *gorm.DB
	Handlers *Handlers
	r        *gin.Engine
	cfg      *config.Config
}

// reads config
func NewApplication() *Application {
	cfg, err := config.Load()
	if err != nil {
		logger.Log.Error("failed to load config", "error", err)
		return nil
	}

	if err := db.Connect(cfg.DSN()); err != nil {
		logger.Log.Error("failed to connect to db", "error", err)
		return nil
	}

	handlers := NewHandlers(db.DB, cfg.JWTSecret)

	r := gin.Default()
	r.LoadHTMLGlob("internal/views/*")
	r.Static("/static", "./static")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	registerAPIroutes(r, handlers, handlers.AuthService)
	registerWebRoutes(r, handlers)

	return &Application{
		DB:       db.DB,
		Handlers: handlers,
		r:        r,
		cfg:      cfg,
	}
}

func (a *Application) Run() {
	if err := a.r.Run(a.cfg.PortAddr()); err != nil {
		logger.Log.Error("failed to run server", "error", err)
	}
}
