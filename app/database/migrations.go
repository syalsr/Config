package database

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() {
	m, err := migrate.New(
		"file://../../schema",
		"postgresql://postgres:postgrespw@localhost:49157/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	
	if err := m.Up(); err != nil {
		log.Println(err)
	}
}
