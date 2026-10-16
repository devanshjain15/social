package main

import (
	"log"
	"time"

	"social/internal/auth"
	"social/internal/db"
	"social/internal/env"
	"social/internal/mailer"
	"social/internal/ratelimiter"
	"social/internal/store"
	"social/internal/store/cache"

	redis "github.com/go-redis/redis/v8"
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
		redis: redisConfig{
			Addr:    env.GetString("REDIS_ADDR", "localhost:6380"),
			Pw:      env.GetString("REDIS_PASSWORD", ""),
			DB:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", true),
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
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
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

	var rdb *redis.Client
	if config.redis.enabled {
		rdb = cache.NewRedisClient(config.redis.Addr, config.redis.Pw, config.redis.DB)
	}

	store := store.NewStorage(db)
	cacheStorage := cache.NewCacheStorage(rdb)

	mailerClient := mailer.NewMailTrapClient(config.mail.mailTrap.apiKey, config.mail.fromEmail)
	authentictor := auth.NewJWTAuthentictor(config.auth.tokenAuth.secret, config.auth.tokenAuth.iss, config.auth.tokenAuth.iss)

	rl := ratelimiter.NewFixedWindowLimiter(config.ratelimiter.RequestsPerTimeFrame, config.ratelimiter.TimeFrame)

	app := &application{
		config:        config,
		store:         store,
		cache:         cacheStorage,
		logger:        logger,
		mailer:        mailerClient,
		authenticator: authentictor,
		ratelimiter: rl,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
