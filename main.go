package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
my_app_up
  - This will be the metric in OCP in Observe -> Metrics
  - Help will be what is displayed as an highlight text.
    ex - my_app_up	prometheus-example-app	web	10.129.0.88:8080	prometheus-example-app	servicemonitor-a	prometheus-example-app-df46974dc-lb2fh	openshift-user-workload-monitoring/user-workload	prometheus-example-app	1
*/
var (
	upGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "my_app_up",
		Help: "1 if the application is up, 0 otherwise",
	})
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "my_app_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(upGauge)
	prometheus.MustRegister(requestCount)
}

func handler(w http.ResponseWriter, r *http.Request) {
	status := "200"
	if r.URL.Path == "/error" {
		status = "500"
	}
	requestCount.WithLabelValues(status).Inc()

	upGauge.Set(1)
}

func main() {
	defer fmt.Println("Application Exiting...")
	fmt.Println("Application... Starting")

	/*
		prom has a handler that can handle /metrics
	*/
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", handler)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			upGauge.Set(0)
			time.Sleep(10 * time.Second)
			upGauge.Set(1)
		}
	}()

	fmt.Println("ListenAndServe... Starting")
	http.ListenAndServe(":8080", nil)
}
