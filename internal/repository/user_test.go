package repository

import (
	"library-api/internal/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {

    dsn := "host=localhost user=postgres password=postgres dbname=library_test port=5433 sslmode=disable"

    db, err := gorm.Open(
        postgres.Open(dsn),
        &gorm.Config{},
    )

    if err != nil {
        t.Fatal(err)
    }

		// refactor
		if err := db.AutoMigrate(
			&models.User{},
			&models.Author{},
			&models.Book{},
			&models.Review{},
			&models.ReadingList{},
		); err != nil {
			t.Fatal("migration failed:", err)
		}

    repo := NewUserRepository(db)

    user := &models.User{
        Name: "John",
        Email: "john@test.com",
    }

    err = repo.Create(user)

    if err != nil {
        t.Fatal(err)
    }

    users, err := repo.GetAll()

    if err != nil {
        t.Fatal(err)
    }

    if len(users) == 0 {
        t.Fatal("expected users")
    }
}
