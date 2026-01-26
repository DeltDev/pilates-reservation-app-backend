package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBSource string
}

func Load() Config {
	_ = godotenv.Load()

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE is not set")
	}

	return Config{
		DBSource: dbSource,
	}
}
