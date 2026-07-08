package main

import (
	"library-api/internal/app"
	"library-api/internal/logger"
	"library-api/migrations"
	"os"
)

func main() {
	a := app.NewApplication()
	if a == nil {
		os.Exit(1)
	}

	if err := migrations.Run(a.DB); err != nil {
		logger.Log.Error("migration failed", "error", err)
		os.Exit(1)
	}

	a.Run()
}
