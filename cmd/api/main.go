package main

import (
	"log"

	"github.com/alanloffler/social/internal/env"
	"github.com/alanloffler/social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
		},
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
