package main

import (
	"log"
	"time"

	"social/internal/auth"
	"social/internal/db"
	"social/internal/env"
	"social/internal/mailer"
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
		addr:   env.GetString("ADDR", ":8000"),
		apiURL: "localhost:8000",
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "PRODUCTION"),
		mail: mailConfig{
			exp:       time.Hour,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basicAuth: basicAuthConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", ""),
				password: env.GetString("BASIC_AUTH_PASSWORD", ""),
			},
			tokenAuth: tokenAuthConfig{
				secret: env.GetString("JWT_SECRET", ""),
				iss:    "gophersocial",
				exp:    time.Hour * 24 * 30,
			},
		},
		frontendURL: env.GetString("FRONTEND_URL", "localhost:3000"),
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

	mailerClient := mailer.NewMailTrapClient(config.mail.mailTrap.apiKey, config.mail.fromEmail)
	authentictor := auth.NewJWTAuthentictor(config.auth.tokenAuth.secret, config.auth.tokenAuth.iss, config.auth.tokenAuth.iss)
	app := &application{
		config:        config,
		store:         store,
		logger:        logger,
		mailer:        mailerClient,
		authenticator: authentictor,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
