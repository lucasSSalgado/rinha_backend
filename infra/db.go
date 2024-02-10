package infra

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateConnection() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=example dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	return db
}
