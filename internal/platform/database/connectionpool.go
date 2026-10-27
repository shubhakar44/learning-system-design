package database

import (
	"context"
	"log"
	"sync"

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

type Pool struct {
	mu               sync.Mutex
	conn             chan pgx.Conn
	maxConns         int
	factory          func() (*pgx.Conn, error)
	totalConnections int
}

func CreateConnectionPool(ctx context.Context, maxConns int) *Pool {
	return &Pool{
		factory: func() (*pgx.Conn, error) {
			return pgx.Connect(ctx, connString)
		},
		conn:             make(chan pgx.Conn, maxConns),
		maxConns:         maxConns,
		totalConnections: 0,
	}
}

func GetPoolInstance(pool *Pool, ctx context.Context) (*pgx.Conn, error) {
	pool.mu.Lock()
	if pool.totalConnections < pool.maxConns {
		conn, err := pool.factory()
		if err != nil {
			pool.totalConnections -= 1
			pool.mu.Unlock()
			return nil, err
		}
		pool.totalConnections += 1
		pool.mu.Unlock()
		return conn, err
	}
	pool.mu.Unlock()
	conn := <-pool.conn
	return &conn, nil
}

func PutPoolInstance(pool *Pool, conn *pgx.Conn) {
	select {
	case pool.conn <- *conn:
	default:
		pool.mu.Lock()
		pool.totalConnections -= 1
		pool.mu.Unlock()
	}
}
