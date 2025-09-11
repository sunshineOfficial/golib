package plugin

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

type Metrics struct {
}

// NewMetrics регистрирует стандартный обработчик Prometheus на роут /metrics.
// По этому роуту можно получить все метрики приложения. Он же будет использоваться
// скрейпером Prometheus для сбора метрик и отправки их на сервер. Все метрики можно
// вывести в Grafana, подключившись к VictoriaMetrics, и на их основе собрать дашборды.
func NewMetrics() *Metrics {
	return &Metrics{}
}

func (p *Metrics) BasePath() string {
	return "/metrics"
}

func (p *Metrics) Register(router *gorouter.Router) {
	router.Handle("", gorouter.WrapStdLib(promhttp.Handler()))
}
