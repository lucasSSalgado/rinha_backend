package infra

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnection() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig("host=postgres port=5432 user=postgres password=example dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	config.MaxConns = int32(7)
	config.MinConns = int32(6)

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	return db
}
