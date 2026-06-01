package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/farrasnazhif/go-std-starter/docs" // this is required to generate swagger docs
	"github.com/farrasnazhif/go-std-starter/internal/mailer"
	"github.com/farrasnazhif/go-std-starter/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config  config
	store   store.Storage
	logger  *zap.SugaredLogger
	mailer  mailer.Client
	metrics *Metrics
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	rateLimiter RateLimitConfig
}

type mailConfig struct {
	resend    resendConfig
	fromEmail string
	exp       time.Duration
}

type resendConfig struct {
	apiKey string
}

type dbConfig struct {
	addr string
	// conns = connections
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string // time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(app.metrics.MetricsMiddleware)

	// Metrics endpoint (exposed outside of /api/v1)
	r.With(app.metricsRateLimiter()).Get("/metrics", MetricsHandler())

	r.Route("/api/v1", func(r chi.Router) {
		// GET /api/v1/health
		r.Get("/health", app.healthCheckHandler)

		// GET /api/v1/swagger
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// GET /api/v1/users/{userID}
		r.Route("/users", func(r chi.Router) {
			// Apply rate limit to activation endpoint
			r.With(app.activateRateLimiter()).Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Use(app.apiRateLimiter())

				r.Get("/", app.getUserHandler)
			})
		})

		// public routes
		r.Route("/auth", func(r chi.Router) {
			// Apply rate limit to registration endpoint
			r.With(app.registerRateLimiter()).Post("/user", app.registerUserHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Version = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/api/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has starter at", "addr", app.config.addr, "env", app.config.env)

	// Create server with graceful shutdown support
	server := NewServer(srv, app)
	return server.ListenAndServeWithGracefulShutdown()
}
