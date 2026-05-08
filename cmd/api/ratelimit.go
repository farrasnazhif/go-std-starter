package main

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/httprate"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	RegisterLimit   int
	ActivateLimit   int
	GeneralAPILimit int
	WindowDuration  time.Duration
}

// DefaultRateLimitConfig returns default rate limit settings
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RegisterLimit:   5,  // 5 requests
		ActivateLimit:   10, // 10 requests
		GeneralAPILimit: 30, // 30 requests
		WindowDuration:  time.Minute,
	}
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Try to get X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// registerRateLimiter returns a rate limiter middleware for registration endpoint
func (app *application) registerRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(
		app.config.rateLimiter.RegisterLimit,
		app.config.rateLimiter.WindowDuration,
	)
}

// activateRateLimiter returns a rate limiter middleware for activation endpoint
func (app *application) activateRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(
		app.config.rateLimiter.ActivateLimit,
		app.config.rateLimiter.WindowDuration,
	)
}

// apiRateLimiter returns a rate limiter middleware for general API endpoints
func (app *application) apiRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(
		app.config.rateLimiter.GeneralAPILimit,
		app.config.rateLimiter.WindowDuration,
	)
}
