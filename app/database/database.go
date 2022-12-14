package database

import (
	"Config/app/config"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Connection struct {
	Pool *pgxpool.Pool
}

var conn *Connection

func GetConnection() *Connection {
	return conn
}

func NewClient(ctx context.Context, sc *config.StorageConfig) (err error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", 
		sc.Username, 
		sc.Password, 
		sc.Host, 
		sc.Port, 
		sc.Database, 
		sc.SSLmode,
	)

	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
		return err
	}
	conn = &Connection{Pool: pool}

	return nil
}
