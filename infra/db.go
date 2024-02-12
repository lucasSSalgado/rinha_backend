package infra

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnection() *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), "host=localhost port=5432 user=postgres password=example dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	return db
}
