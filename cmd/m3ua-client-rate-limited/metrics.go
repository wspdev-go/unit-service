package main

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ClientReceiveTotal prometheus.Counter
var ClientSendTotal prometheus.Counter
var ClientReceiveTotalVec *prometheus.CounterVec
var ClientSendTotalVec *prometheus.CounterVec

const clientNamespace = "hstp"
const clientSubsystem = "client"

func init() {
	ClientReceiveTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: clientNamespace,
		Subsystem: clientSubsystem,
		Name:      "recv_total",
		Help:      "The total number of received messages",
	})

	ClientSendTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: clientNamespace,
		Subsystem: clientSubsystem,
		Name:      "send_total",
		Help:      "The total number of sended messages",
	})

	ClientReceiveTotalVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: clientNamespace,
			Subsystem: clientSubsystem,
			Name:      "recv_total_vec",
			Help:      "The total number of received messages",
		},
		[]string{"hstp_message"})

	ClientSendTotalVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: clientNamespace,
			Subsystem: clientSubsystem,
			Name:      "send_total_vec",
			Help:      "The total number of sended messagess",
		},
		[]string{"hstp_message"})

	go run()
}

func run() {

	clientAddr := os.Getenv("CLIENT_PROMETHEUS_PORT")
	if clientAddr == "" {
		clientAddr = ":11001"
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(ClientReceiveTotal)
	reg.MustRegister(ClientSendTotal)
	reg.MustRegister(ClientReceiveTotalVec)
	reg.MustRegister(ClientSendTotalVec)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

}
