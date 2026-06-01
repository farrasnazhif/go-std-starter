package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	httpRequestDuration prometheus.Histogram
	httpRequestCount    prometheus.Counter
	httpRequestsActive  prometheus.Gauge
	httpStatusCodes     prometheus.CounterVec
	dbConnections       prometheus.Gauge
	dbQueriesDuration   prometheus.Histogram
}

func NewMetrics() (*Metrics, error) {
	m := &Metrics{
		httpRequestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		httpRequestCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		}),
		httpRequestsActive: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Current number of active HTTP requests",
		}),
		httpStatusCodes: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_responses_total",
			Help: "Total number of HTTP responses by status code",
		}, []string{"status_code", "method", "path"}),
		dbConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Current number of active database connections",
		}),
		dbQueriesDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: prometheus.DefBuckets,
		}),
	}

	collectors := []prometheus.Collector{
		m.httpRequestDuration, m.httpRequestCount, m.httpRequestsActive,
		&m.httpStatusCodes, m.dbConnections, m.dbQueriesDuration,
	}
	for _, c := range collectors {
		if err := prometheus.Register(c); err != nil {
			return nil, fmt.Errorf("failed to register metric: %w", err)
		}
	}

	return m, nil
}

func (m *Metrics) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		m.httpRequestsActive.Inc()
		defer m.httpRequestsActive.Dec()

		m.httpRequestCount.Inc()

		start := time.Now()
		defer func() {
			m.httpRequestDuration.Observe(time.Since(start).Seconds())
		}()

		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(sw, r)

		m.httpStatusCodes.WithLabelValues(strconv.Itoa(sw.statusCode), r.Method, r.URL.Path).Inc()
	})
}

func (m *Metrics) RecordDBConnection(count int) {
	m.dbConnections.Set(float64(count))
}

func (m *Metrics) RecordDBQueryDuration(duration time.Duration) {
	m.dbQueriesDuration.Observe(duration.Seconds())
}

func MetricsHandler() http.HandlerFunc {
	h := promhttp.Handler()
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (w *statusWriter) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.written = true
	}
	return w.ResponseWriter.Write(b)
}
