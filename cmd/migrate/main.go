package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up', 'down', or 'force [version]'")
	}

	direction := os.Args[1]

	// Connect to the default 'postgres' database to check if 'makerble' exists
	initialDsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := sql.Open("postgres", initialDsn)
	if err != nil {
		log.Fatal("Failed to connect to postgres database:", err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'makerble')").Scan(&exists)
	if err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	if !exists {
		_, err = db.Exec("CREATE DATABASE makerble")
		if err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Println("Database 'makerble' created successfully")
	} else {
		log.Println("Database 'makerble' already exists")
	}

	// Now, connect to the 'makerble' database to run migrations
	dsn := "host=localhost user=postgres password=postgres dbname=makerble port=5432 sslmode=disable"
	migrateDb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to makerble database for migration:", err)
	}
	defer migrateDb.Close()

	instance, err := postgres.WithInstance(migrateDb, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	default:
		log.Fatal("Invalid direction")
	}

	log.Println("Migration completed successfully.")
}
