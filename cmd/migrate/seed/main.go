package main

import (
	"log"

	"github.com/alanloffler/social/internal/db"
	"github.com/alanloffler/social/internal/env"
	"github.com/alanloffler/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://postgres:admin1234@localhost/social?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)
}
