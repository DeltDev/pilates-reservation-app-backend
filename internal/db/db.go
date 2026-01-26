package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dbSource string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	return pool
}
