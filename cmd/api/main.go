package main

import (
	"flag"
	"library-api/internal/app"
	"library-api/internal/logger"
	"library-api/migrations"
	"os"
)

func main() {
	migrateOnly := flag.Bool("migrate-only", false, "run pending database migrations and exit, without starting the API server")
	flag.Parse()

	a := app.NewApplication()
	if a == nil {
		os.Exit(1)
	}

	if *migrateOnly {
		if err := migrations.Run(a.DB); err != nil {
			logger.Log.Error("migration failed", "error", err)
			os.Exit(1)
		}
		return
	}

	a.Run()
}
