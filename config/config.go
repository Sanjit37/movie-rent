package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func LoadDBConfig() DBConfig {
	_ = godotenv.Load() // Load .env if exists

	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func (c DBConfig) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}
