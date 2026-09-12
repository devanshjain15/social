package main

import (
	"log"

	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

const VERSION string = "0.0.1"

func main() {
	config := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:secret@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "DEV"),
	}

	db, err := db.New(config.db.addr, config.db.maxOpenConns, config.db.maxIdleConns, config.db.maxIdleTime)
	if err != nil {
		log.Fatalf("Database unreachable: %v", err)
	} else {
		log.Println("Database connection pool established!")
	}
	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: config,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
