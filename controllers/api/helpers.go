package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/heyjoakim/devops-21/services"
	"github.com/prometheus/client_golang/prometheus"
)

func updateLatest(r *http.Request) {
	tryLatestQuery := r.URL.Query().Get("latest")

	if tryLatestQuery != "" {
		tryLatest, _ := strconv.Atoi(tryLatestQuery)

		services.UpdateLatest(tryLatest)
	}
}

func RegisterEndpoint(name string) *prometheus.Timer {
	hist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    fmt.Sprintf("http_request_%s_duration_seconds", name),
			Help:    fmt.Sprintf("http_request_%s_duration_seconds", name),
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
		},
		[]string{"status"},
	)
	var status string
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		hist.WithLabelValues(status).Observe(v)
	}))
	prometheus.MustRegister(hist)
	return timer
}
