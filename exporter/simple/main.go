package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var constLabel = prometheus.Labels{"service": "manager"}

var vLable = []string{"x_service"}

var vValue = []string{"x_alive"}

var status = 0

type StatusCollector struct {
	statusMetric *prometheus.Desc
}

func (collector *StatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.statusMetric
}

func (collector *StatusCollector) Collect(ch chan<- prometheus.Metric) {
	metricValue := float64(status)
	ch <- prometheus.MustNewConstMetric(collector.statusMetric, prometheus.CounterValue, metricValue, vValue...)
}

func NewStatusCollector() *StatusCollector {

	return &StatusCollector{
		// metric name/metric label values/metric help text/metric type/measurement
		statusMetric: prometheus.NewDesc("status_metric", "Shows status of the service", vLable, constLabel),
	}
}

func changeStatus() {
	rand.Seed(time.Now().UnixNano())
	for {
		n := rand.Int() % 10
		time.Sleep(time.Duration(n) * time.Second)
		if status == 0 {
			status = 1
		} else {
			status = 0
		}
	}
}

func main() {
	sc := NewStatusCollector()
	prometheus.MustRegister(sc)

	go changeStatus()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9999", nil))
}
