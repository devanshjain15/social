package main

import (
	"log"
	"time"

	"social/internal/db"
	"social/internal/env"
	"social/internal/store"

	"go.uber.org/zap"
)

const VERSION string = "0.0.1"

//	@title			GopherSocial
//	@description	API for the GopherSocial, a social media for gophers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	config := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: "localhost:8000",
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:secret@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "DEV"),
		mail: mailConfig{
			exp: time.Hour,
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

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
		logger: logger,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
