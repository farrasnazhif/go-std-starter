package main

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	RegisterLimit   int
	ActivateLimit   int
	GeneralAPILimit int
	MetricsLimit    int
	WindowDuration  time.Duration
}

// DefaultRateLimitConfig returns default rate limit settings
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RegisterLimit:   5,
		ActivateLimit:   10,
		GeneralAPILimit: 30,
		MetricsLimit:    15,
		WindowDuration:  time.Minute,
	}
}

func (app *application) rateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		limit,
		window,
		httprate.WithKeyFuncs(httprate.KeyByRealIP),
		httprate.WithLimitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.rateLimitExceededResponse(w, r)
		})),
	)
}

func (app *application) registerRateLimiter() func(http.Handler) http.Handler {
	return app.rateLimiter(app.config.rateLimiter.RegisterLimit, app.config.rateLimiter.WindowDuration)
}

func (app *application) activateRateLimiter() func(http.Handler) http.Handler {
	return app.rateLimiter(app.config.rateLimiter.ActivateLimit, app.config.rateLimiter.WindowDuration)
}

func (app *application) apiRateLimiter() func(http.Handler) http.Handler {
	return app.rateLimiter(app.config.rateLimiter.GeneralAPILimit, app.config.rateLimiter.WindowDuration)
}

func (app *application) metricsRateLimiter() func(http.Handler) http.Handler {
	return app.rateLimiter(app.config.rateLimiter.MetricsLimit, app.config.rateLimiter.WindowDuration)
}
