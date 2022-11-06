package database

import (
	"Config/app/config"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(sc *config.StorageConfig) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		sc.Username,
		sc.Password,
		sc.Host,
		sc.Port,
		sc.Database,
		sc.SSLmode,
	)
	fileMigrations := "file://schema"

	m, err := migrate.New(fileMigrations, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Println(err)
	}
}
