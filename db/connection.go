package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"movie-rent/config"
)

type Database interface {
	Instance() *sqlx.DB
}

type database struct{}

func NewDatabase() Database {
	return database{}
}

func (d database) Instance() *sqlx.DB {
	connStr := config.LoadDBConfig().GetDSN()
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
