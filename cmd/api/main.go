package main

import (
	"github.com/farrasnazhif/go-std-starter/cmd/api/middleware"
	"github.com/farrasnazhif/go-std-starter/internal/db"
	"github.com/farrasnazhif/go-std-starter/internal/env"
	"github.com/farrasnazhif/go-std-starter/internal/mailer"
	"github.com/farrasnazhif/go-std-starter/internal/service"
	"github.com/farrasnazhif/go-std-starter/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			SocialNetwork API
//	@description	API for SocialNetwork
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/api/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/go-std-starter?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			fromEmail: env.GetString("MAIL_FROM_EMAIL", "onboarding@resend.dev"),
			resend: resendConfig{
				apiKey: env.GetString("RESEND_API_KEY", ""),
			},
		},
		rateLimiter: middleware.DefaultRateLimitConfig(),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	database, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer database.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(database)
	mailer := mailer.NewResend(cfg.mail.resend.apiKey, cfg.mail.fromEmail)

	metrics, err := middleware.NewMetrics()
	if err != nil {
		logger.Fatal(err)
	}

	authService := service.NewAuthService(store, mailer, cfg.env)
	otpService := service.NewOTPService(store, mailer, cfg.env)

	app := &application{
		config:      cfg,
		store:       store,
		logger:      logger,
		mailer:      mailer,
		metrics:     metrics,
		authService: authService,
		otpService:  otpService,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
