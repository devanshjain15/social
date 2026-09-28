package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "")
	maxOpenConns := env.GetInt("DB_MAX_OPEN_CONNS", 30)
	maxIdleConns := env.GetInt("DB_MAX_OPEN_CONNS", 30)
	maxIdleTime := env.GetString("DB_MAX_IDLE_TIME", "15m")
	db, err := db.New(addr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	defer db.Close()
	store := store.NewStorage(db)
	Seed(store, db)
}
