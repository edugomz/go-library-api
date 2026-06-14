package main

import (
	"library-api/internal/app"
	"library-api/internal/logger"
	"library-api/migrations"
)

func main() {

	app := app.NewApplication()

	if err := migrations.Run(app.DB); err != nil {
		logger.Log.Error("migration failed", "error", err)
	}

	app.Run()
}
