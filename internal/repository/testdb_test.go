package repository

import (
	"sync/atomic"
	"testing"
	"time"

	"library-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var uniqueCounter int64

// migrateLockKey serializes AutoMigrate across concurrently running test
// binaries (this package and internal/app's e2e tests both migrate the same
// shared DB, and go test ./... runs separate packages' binaries in
// parallel). Must match the key used in internal/app/e2e_test.go.
const migrateLockKey = 918273645

// uniqueSuffix returns a value unique within this test run, used to build
// unique emails/names since the test DB isn't truncated between runs.
func uniqueSuffix() int64 {
	return time.Now().UnixNano() + atomic.AddInt64(&uniqueCounter, 1)
}

// setupTestDB opens a connection to the docker-compose.test.yml Postgres
// instance and ensures the schema is migrated. Tests share this DB and do
// not truncate between runs, so assertions should check for presence/values
// of records they created rather than exact row counts.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "host=localhost user=postgres password=postgres dbname=library_test port=5433 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", migrateLockKey).Error; err != nil {
			return err
		}
		return tx.AutoMigrate(
			&models.User{},
			&models.Author{},
			&models.Book{},
			&models.Review{},
			&models.ReadingList{},
		)
	}); err != nil {
		t.Fatal("migration failed:", err)
	}

	return db
}
