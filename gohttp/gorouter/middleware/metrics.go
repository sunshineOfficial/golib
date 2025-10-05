package middleware

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

type PrometheusMetric interface {
	Collector() prometheus.Collector
}

type techMetrics struct {
	totalRequests   *prometheus.CounterVec
	activeRequests  *prometheus.GaugeVec
	requestDuration *prometheus.HistogramVec
}

// Metrics включает базовые технические метрики и дает возможность зарегистрировать кастомные метрики.
// Функцию можно вызвать только один раз при инициализации приложения, иначе будет паника.
func Metrics(metrics ...PrometheusMetric) gorouter.Middleware {
	if len(metrics) != 0 {
		cs := make([]prometheus.Collector, 0, len(metrics))
		for _, metric := range metrics {
			cs = append(cs, metric.Collector())
		}

		prometheus.MustRegister(cs...)
	}

	m := techMetrics{
		totalRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"path", "status"},
		),
		activeRequests: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_active_requests",
				Help: "Number of currently active HTTP requests",
			},
			[]string{"path"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Histogram of HTTP request durations in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 30, 60},
			},
			[]string{"path"},
		),
	}

	return func(_ gorouter.Context, h gorouter.Handler) gorouter.Handler {
		return func(c gorouter.Context) error {
			var (
				start = time.Now()
				path  = c.PathTemplate()
			)

			m.activeRequests.WithLabelValues(path).Inc()

			defer func() {
				m.requestDuration.WithLabelValues(path).Observe(time.Since(start).Seconds())

				status := c.Response().Status()

				m.totalRequests.WithLabelValues(path, strconv.Itoa(status)).Inc()
				m.activeRequests.WithLabelValues(path).Dec()
			}()

			return h(c)
		}
	}
}
