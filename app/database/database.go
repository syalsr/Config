package database

import (
	"Config/app/config"
	"context"
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
	//dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	dsn := "postgresql://postgres:postgrespw@localhost:49157/postgres?sslmode=disable"
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
		return err
	}
	conn = &Connection{Pool: pool}

	return nil
}
