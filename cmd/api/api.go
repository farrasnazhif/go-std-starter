package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/farrasnazhif/go-std-starter/cmd/api/middleware"
	"github.com/farrasnazhif/go-std-starter/docs"
	"github.com/farrasnazhif/go-std-starter/internal/mailer"
	"github.com/farrasnazhif/go-std-starter/internal/service"
	"github.com/farrasnazhif/go-std-starter/internal/store"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config      config
	store       store.Storage
	logger      *zap.SugaredLogger
	mailer      mailer.Client
	metrics     *middleware.Metrics
	authService *service.AuthService
	otpService  *service.OTPService
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	rateLimiter middleware.RateLimitConfig
}

type mailConfig struct {
	resend    resendConfig
	fromEmail string
}

type resendConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Recoverer)
	r.Use(chimw.Logger)
	r.Use(middleware.CORS(middleware.DefaultCORSConfig(app.config.frontendURL)))
	r.Use(app.metrics.Handler)

	r.With(app.rateLimiter(app.config.rateLimiter.MetricsLimit)).Get("/metrics", middleware.MetricsHandler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Use(app.rateLimiter(app.config.rateLimiter.GeneralAPILimit))
				r.Get("/", app.getUserHandler)
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.With(app.rateLimiter(app.config.rateLimiter.RegisterLimit)).Post("/user", app.registerUserHandler)
			r.Post("/forgot-password", app.forgotPasswordHandler)
			r.Post("/forgot-password/verify", app.verifyForgotPasswordHandler)
			r.Post("/reset-password", app.resetPasswordHandler)
		})

		r.Route("/otp", func(r chi.Router) {
			r.With(app.rateLimiter(app.config.rateLimiter.OTPSendLimit)).Post("/send", app.sendOTPHandler)
			r.Post("/verify", app.verifyOTPHandler)
		})
	})

	return r
}

func (app *application) rateLimiter(limit int) func(http.Handler) http.Handler {
	return middleware.RateLimiter(limit, app.config.rateLimiter.WindowDuration, app.rateLimitExceededResponse)
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/api/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has started at", "addr", app.config.addr, "env", app.config.env)

	server := NewServer(srv, app)
	return server.ListenAndServeWithGracefulShutdown()
}
