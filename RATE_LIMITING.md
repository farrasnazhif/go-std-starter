# Rate Limiting Configuration

The API includes built-in rate limiting to protect against abuse and brute force attacks.

## Default Limits

- **Registration (`POST /api/v1/auth/user`)**: 5 requests per minute per IP
- **Activation (`PUT /api/v1/users/activate/{token}`)**: 10 requests per minute per IP
- **General API endpoints**: 30 requests per minute per IP

## Configuration

Rate limits are configured in `cmd/api/ratelimit.go` and can be customized by modifying the `DefaultRateLimitConfig()` function:

```go
RateLimitConfig{
    RegisterLimit:    5,             // requests per window
    ActivateLimit:    10,            // requests per window
    GeneralAPILimit:  30,            // requests per window
    WindowDuration:   time.Minute,   // time window
}
```

## How It Works

- Rate limiting is based on **client IP address**
- Supports proxied requests via `X-Forwarded-For` and `X-Real-IP` headers
- Uses in-memory token bucket algorithm
- Returns `429 Too Many Requests` when limit exceeded

## Customization

To customize rate limits, modify the `DefaultRateLimitConfig()` function in `ratelimit.go`:

```go
func DefaultRateLimitConfig() RateLimitConfig {
    return RateLimitConfig{
        RegisterLimit:   10,  // Increase from 5 to 10
        ActivateLimit:   20,  // Increase from 10 to 20
        GeneralAPILimit: 60,  // Increase from 30 to 60
        WindowDuration:  time.Minute,
    }
}
```

## Adding Rate Limits to New Endpoints

To add rate limiting to a new endpoint:

```go
r.Route("/new-endpoint", func(r chi.Router) {
    r.With(app.apiRateLimiter()).Post("/", app.newEndpointHandler)
})
```

Or create a custom rate limiter:

```go
customLimiter := httprate.LimitByIP(20, time.Minute)
r.With(customLimiter).Post("/custom", app.customHandler)
```
