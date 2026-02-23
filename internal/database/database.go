package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)
var DB *pgxpool.Pool

func Connect(dbUrl string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("DB Ping failed: %v", err)
	}

	DB = pool
	log.Println("Database connected successfully")
}