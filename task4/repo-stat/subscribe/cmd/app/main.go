package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	//config ...

	/*
			dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		    cfg.Database.User, cfg.Database.Password, cfg.Database.Host,
		    cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)
	*/
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://user:pass@localhost:5432/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
