package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define custom metrics
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
	// Register custom metrics
	prometheus.MustRegister(upGauge)
	prometheus.MustRegister(requestCount)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Increment request counter based on the status
	status := "200"
	if r.URL.Path == "/error" {
		status = "500"
	}
	requestCount.WithLabelValues(status).Inc()

	// Set the 'up' gauge metric
	upGauge.Set(1)

	// Respond with a simple message
	fmt.Fprintf(w, "Hello from the custom Go app!\n")
}

func main() {
	// Expose the /metrics endpoint for Prometheus scraping
	http.Handle("/metrics", promhttp.Handler())

	// Define HTTP routes for your application
	http.HandleFunc("/", handler)

	// Start HTTP server
	go func() {
		for {
			// Simulate the app being down periodically
			time.Sleep(10 * time.Second)
			upGauge.Set(0)
			time.Sleep(10 * time.Second)
			upGauge.Set(1)
		}
	}()

	fmt.Println("ListenAndServe... Starting")
	// Start the server
	http.ListenAndServe(":8080", nil)
}
