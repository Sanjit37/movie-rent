package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type Database interface {
	Instance() *sqlx.DB
}

type database struct{}

func NewDatabase() Database {
	return database{}
}

var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Sanjit"
	dbname   = "movie_db"
)

func (d database) Instance() *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := CheckDBConnection(db); err != nil {
		log.Fatal(err)
	}

	return db
}

func CheckDBConnection(db *sqlx.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("❌ database is not connected: %w", err)
	}
	fmt.Println("✅ database is connected")
	return nil
}
