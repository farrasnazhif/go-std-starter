package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

type RateLimitConfig struct {
	RegisterLimit   int
	GeneralAPILimit int
	MetricsLimit    int
	OTPSendLimit    int
	WindowDuration  time.Duration
}

func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RegisterLimit:   5,
		GeneralAPILimit: 30,
		MetricsLimit:    15,
		OTPSendLimit:    3,
		WindowDuration:  time.Minute,
	}
}

// RateLimitExceededHandler is the function signature for handling 429 responses.
type RateLimitExceededHandler func(w http.ResponseWriter, r *http.Request)

func RateLimiter(limit int, window time.Duration, onLimitExceeded RateLimitExceededHandler) func(http.Handler) http.Handler {
	return httprate.Limit(
		limit,
		window,
		httprate.WithKeyFuncs(httprate.KeyByRealIP),
		httprate.WithLimitHandler(http.HandlerFunc(onLimitExceeded)),
	)
}
