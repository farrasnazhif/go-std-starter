# Server Metrics

The API collects and exposes Prometheus metrics for monitoring and observability.

## Metrics Endpoint

Metrics are exposed at `/metrics` in Prometheus format:

```
GET /metrics
```

## Available Metrics

### HTTP Request Metrics

#### `http_requests_total` (Counter)

Total number of HTTP requests processed.

**Labels:**

- No additional labels

**Example:**

```
http_requests_total 1250
```

#### `http_request_duration_seconds` (Histogram)

HTTP request latency in seconds.

**Buckets:** .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10

**Example:**

```
http_request_duration_seconds_bucket{le="0.05"} 456
http_request_duration_seconds_bucket{le="0.1"} 789
http_request_duration_seconds_bucket{le="+Inf"} 1250
http_request_duration_seconds_sum 234.567
http_request_duration_seconds_count 1250
```

#### `http_requests_active` (Gauge)

Current number of active HTTP requests being processed.

**Example:**

```
http_requests_active 12
```

#### `http_responses_total` (Counter)

Total number of HTTP responses by status code, method, and path.

**Labels:**

- `status_code`: HTTP status code (200, 400, 401, 429, 500, etc.)
- `method`: HTTP method (GET, POST, PUT, DELETE)
- `path`: Request path

**Example:**

```
http_responses_total{method="POST",path="/api/v1/auth/user",status_code="201"} 145
http_responses_total{method="POST",path="/api/v1/auth/user",status_code="400"} 23
http_responses_total{method="GET",path="/api/v1/users/1",status_code="200"} 567
```

### Database Metrics

#### `db_connections_active` (Gauge)

Current number of active database connections.

**Example:**

```
db_connections_active 8
```

#### `db_query_duration_seconds` (Histogram)

Database query latency in seconds.

**Buckets:** .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10

**Example:**

```
db_query_duration_seconds_bucket{le="0.01"} 234
db_query_duration_seconds_bucket{le="0.05"} 512
db_query_duration_seconds_bucket{le="+Inf"} 789
db_query_duration_seconds_sum 45.123
db_query_duration_seconds_count 789
```

## Monitoring with Prometheus

### Basic Prometheus Configuration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "go-std-starter"
    static_configs:
      - targets: ["localhost:8080"]
    metrics_path: "/metrics"
    scrape_interval: 15s
```

### Useful Queries

#### Request Rate (requests per second)

```promql
rate(http_requests_total[1m])
```

#### P95 Request Latency

```promql
histogram_quantile(0.95, http_request_duration_seconds_bucket)
```

#### Error Rate

```promql
rate(http_responses_total{status_code=~"5.."}[1m])
```

#### Active Requests

```promql
http_requests_active
```

#### Average Response Time

```promql
rate(http_request_duration_seconds_sum[1m]) / rate(http_request_duration_seconds_count[1m])
```

#### Status Code Distribution

```promql
http_responses_total
```

#### Registration Endpoint Success Rate

```promql
rate(http_responses_total{path="/api/v1/auth/user",status_code="201"}[5m]) / ignoring(status_code) group_left sum by (path) (rate(http_responses_total{path="/api/v1/auth/user"}[5m]))
```

## Grafana Integration

### Sample Dashboard Variables

```json
{
  "datasource": "Prometheus",
  "targets": [
    {
      "expr": "rate(http_requests_total[1m])",
      "legendFormat": "Requests/sec"
    },
    {
      "expr": "histogram_quantile(0.95, http_request_duration_seconds_bucket)",
      "legendFormat": "P95 Latency"
    },
    {
      "expr": "http_requests_active",
      "legendFormat": "Active Requests"
    }
  ]
}
```

## Best Practices

1. **Scrape Interval**: Set to 15-30 seconds for balanced resource usage and data granularity
2. **Retention**: Keep at least 2 weeks of metrics data for trend analysis
3. **Alerting**: Set up alerts for:
   - High error rates (5xx responses)
   - High request latency (P95 > 500ms)
   - High rate limit violations (429 responses)
   - Database connection pool exhaustion

4. **Cardinality**: Monitor label cardinality to prevent high memory usage:
   - Use appropriate scrape intervals
   - Limit dimensions on custom metrics
   - Set up recording rules for aggregation

## Example Alert Rules

```yaml
groups:
  - name: go-std-starter
    rules:
      - alert: HighErrorRate
        expr: |
          rate(http_responses_total{status_code=~"5.."}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "High error rate detected"

      - alert: HighLatency
        expr: |
          histogram_quantile(0.95, http_request_duration_seconds_bucket) > 0.5
        for: 5m
        annotations:
          summary: "High request latency detected"

      - alert: RateLimitExceeded
        expr: |
          rate(http_responses_total{status_code="429"}[5m]) > 0.1
        for: 2m
        annotations:
          summary: "Rate limit violations detected"
```

## Disabling Metrics

To disable metrics collection, comment out or remove the metrics middleware from `cmd/api/api.go`:

```go
// r.Use(app.metrics.MetricsMiddleware)
```

And remove the metrics endpoint:

```go
// r.Get("/metrics", MetricsHandler())
```
