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
	requestDuration *prometheus.SummaryVec
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
		requestDuration: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_request_duration_summary",
				Help: "Summary of duration of HTTP requests in seconds",
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
