package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type BrokerMetrics struct {
	RpcMethodCount    *prometheus.CounterVec
	RpcMethodLatency  *prometheus.HistogramVec
	ActiveSubscribers prometheus.Gauge
}

func StartPrometheus() BrokerMetrics {
	rpcMethodCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "broker_rpc_method_count",
			Help: "Count of rpc method calls",
		},
		[]string{"method", "status"},
	)

	rpcMethodLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "broker_rpc_method_latency",
			Help:    "Latency of rpc methods",
			Buckets: []float64{0.125, 0.25, 0.5, 0.75, 1, 1.25, 1.5, 2},
		},
		[]string{"method"},
	)
	activeSubscribers := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "broker_active_subscribers",
			Help: "Number of active subscribers",
		},
	)
	metrics := BrokerMetrics{
		RpcMethodCount:    rpcMethodCounter,
		RpcMethodLatency:  rpcMethodLatency,
		ActiveSubscribers: activeSubscribers,
	}
	prometheus.MustRegister(rpcMethodCounter)
	prometheus.MustRegister(rpcMethodLatency)
	prometheus.MustRegister(activeSubscribers)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":9091", nil)
		if err != nil {
			log.Errorln(err)
			return
		}
	}()
	return metrics
}
