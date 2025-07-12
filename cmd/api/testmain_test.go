package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"makerble_api/internal/database"
	"makerble_api/internal/env"

	_ "github.com/lib/pq"
)

var testApp application

func TestMain(m *testing.M) {
	dsn := "host=localhost user=postgres password=postgres dbname=makerble_test port=5432 sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	testApp = application{
		port:      env.GetEnvInt("PORT", 8081),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    database.NewModels(db),
		DB:        db,
	}

	code := m.Run()

	_, err = db.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatal("Failed to clean up test database:", err)
	}

	os.Exit(code)
}
