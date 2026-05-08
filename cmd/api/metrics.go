package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all application metrics
type Metrics struct {
	httpRequestDuration prometheus.Histogram
	httpRequestCount    prometheus.Counter
	httpRequestsActive  prometheus.Gauge
	httpStatusCodes     prometheus.CounterVec
	dbConnections       prometheus.Gauge
	dbQueriesDuration   prometheus.Histogram
}

// NewMetrics initializes and registers Prometheus metrics
func NewMetrics() (*Metrics, error) {
	m := &Metrics{
		httpRequestDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
		),
		httpRequestCount: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
		),
		httpRequestsActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_active",
				Help: "Current number of active HTTP requests",
			},
		),
		httpStatusCodes: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_responses_total",
				Help: "Total number of HTTP responses by status code",
			},
			[]string{"status_code", "method", "path"},
		),
		dbConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Current number of active database connections",
			},
		),
		dbQueriesDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
		),
	}

	// Register all metrics
	if err := prometheus.Register(m.httpRequestDuration); err != nil {
		return nil, fmt.Errorf("failed to register httpRequestDuration: %w", err)
	}
	if err := prometheus.Register(m.httpRequestCount); err != nil {
		return nil, fmt.Errorf("failed to register httpRequestCount: %w", err)
	}
	if err := prometheus.Register(m.httpRequestsActive); err != nil {
		return nil, fmt.Errorf("failed to register httpRequestsActive: %w", err)
	}
	if err := prometheus.Register(&m.httpStatusCodes); err != nil {
		return nil, fmt.Errorf("failed to register httpStatusCodes: %w", err)
	}
	if err := prometheus.Register(m.dbConnections); err != nil {
		return nil, fmt.Errorf("failed to register dbConnections: %w", err)
	}
	if err := prometheus.Register(m.dbQueriesDuration); err != nil {
		return nil, fmt.Errorf("failed to register dbQueriesDuration: %w", err)
	}

	return m, nil
}

// MetricsMiddleware returns an HTTP middleware that records request metrics
func (m *Metrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip metrics endpoint itself
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		// Increment active requests
		m.httpRequestsActive.Inc()
		defer m.httpRequestsActive.Dec()

		// Increment total request count
		m.httpRequestCount.Inc()

		// Record request duration
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			m.httpRequestDuration.Observe(duration)
		}()

		// Wrap response writer to capture status code
		wrappedWriter := &statusCodeWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		// Record status code
		statusCode := strconv.Itoa(wrappedWriter.statusCode)
		m.httpStatusCodes.WithLabelValues(statusCode, r.Method, r.URL.Path).Inc()
	})
}

// statusCodeWriter wraps http.ResponseWriter to capture the status code
type statusCodeWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (w *statusCodeWriter) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *statusCodeWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.written = true
	}
	return w.ResponseWriter.Write(b)
}

// RecordDBConnection updates the active database connections count
func (m *Metrics) RecordDBConnection(count int) {
	m.dbConnections.Set(float64(count))
}

// RecordDBQueryDuration records the duration of a database query
func (m *Metrics) RecordDBQueryDuration(duration time.Duration) {
	m.dbQueriesDuration.Observe(duration.Seconds())
}

// Handler returns an HTTP handler function for the /metrics endpoint
func MetricsHandler() http.HandlerFunc {
	metricsHandler := promhttp.Handler()
	return func(w http.ResponseWriter, r *http.Request) {
		metricsHandler.ServeHTTP(w, r)
	}
}
