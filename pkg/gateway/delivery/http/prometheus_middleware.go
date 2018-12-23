package http

import (
	"runtime"
	"strconv"
	"time"

	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/wakanda/pkg/gateway"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMiddleware struct {
	manager *gateway.Manager

	// prometheus funcs
	bootGaugeFunc                prometheus.GaugeFunc
	upTimeGaugeFunc              prometheus.GaugeFunc
	cpuCountGaugeFunc            prometheus.GaugeFunc
	onlinePeopleFunc             prometheus.GaugeFunc
	httpRequestCounterCounterVec *prometheus.CounterVec
	requestDurationHistogram     prometheus.Histogram
	startAt                      time.Time
}

func NewPrometheusMiddleware(manager *gateway.Manager) *PrometheusMiddleware {
	m := PrometheusMiddleware{
		manager: manager,
	}
	m.startAt = time.Now().UTC()

	m.bootGaugeFunc = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "boot_time",
			Help: "boot time in unix type",
		},
		func() float64 { return float64(m.startAt.Unix()) },
	)
	prometheus.MustRegister(m.bootGaugeFunc)

	m.onlinePeopleFunc = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "online_people",
			Help: "online people",
		},
		func() float64 { return float64(m.manager.Status().OnlinePeople) },
	)
	prometheus.MustRegister(m.onlinePeopleFunc)

	m.upTimeGaugeFunc = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "uptime",
			Help: "elapse time since boot in unix type",
		},
		func() float64 { return float64(time.Now().UTC().Unix() - m.startAt.Unix()) },
	)
	prometheus.MustRegister(m.upTimeGaugeFunc)

	m.cpuCountGaugeFunc = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "cpu_core",
			Help: "number of cpu",
		},
		func() float64 { return float64(runtime.NumCPU()) },
	)
	prometheus.MustRegister(m.cpuCountGaugeFunc)

	m.httpRequestCounterCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request",
			Help: "number of http requests group by status",
		},
		[]string{"status"},
	)
	prometheus.MustRegister(m.httpRequestCounterCounterVec)

	m.requestDurationHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "time of duration in second type",
		},
	)
	prometheus.MustRegister(m.requestDurationHistogram)

	return &m
}

func (m *PrometheusMiddleware) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	startTime := time.Now()
	next(c)
	duration := float64(time.Since(startTime)/time.Millisecond) / 1000
	m.requestDurationHistogram.Observe(duration)
	m.httpRequestCounterCounterVec.WithLabelValues(strconv.Itoa(c.Status())).Inc()
}

func NewPrometheusRouter() *napnap.Router {
	router := napnap.NewRouter()
	router.Get("/metrics", napnap.WrapHandler(promhttp.Handler()))
	return router
}
