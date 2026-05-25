package metrics

import (
	"net/http"
	"unit-service/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "transaction"
	subsystem = "trace"
	addr      = "localhost:11001"
)

var TransactionInVec *prometheus.CounterVec

func init() {
	TransactionInVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "in_vec",
			Help:      "The vec number of Transaction In",
		},
		[]string{"transaction_message"})
	go run()
}

func run() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(TransactionInVec)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("Unable to start server: %v", err)
	}
}
