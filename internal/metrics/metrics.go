package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "code"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	RateLimitedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limited_total",
			Help: "Total requests that were rate limited (429)",
		},
		[]string{"path"},
	)

	RedisOpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "redis_op_duration_seconds",
			Help: "Redis operation duration in seconds",
			Buckets: []float64{0.0005, 0.001, 0.002, 0.005, 0.01, 0.02, 0.05},
		},
		[]string{"op"},
	)
)