package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBSource string
	MidtransServerKey  string
	MidtransClientKey  string
	MidtransEnvironment string
}

func Load() Config {
	_ = godotenv.Load()

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE is not set")
	}

	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if midtransServerKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY is not set")
	}

	midtransClientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	if midtransClientKey == "" {
		log.Fatal("MIDTRANS_CLIENT_KEY is not set")
	}

	midtransEnv := os.Getenv("MIDTRANS_ENVIRONMENT")
	if midtransEnv == "" {
		midtransEnv = "sandbox" 
	}

	return Config{
		DBSource:            dbSource,
		MidtransServerKey:   midtransServerKey,
		MidtransClientKey:   midtransClientKey,
		MidtransEnvironment: midtransEnv,
	}
}
