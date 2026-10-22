package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var connString = "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable"

func CreateInstance(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
