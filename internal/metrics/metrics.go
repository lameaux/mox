package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const (
	readHeaderTimeout = 5 * time.Second

	path = "/metrics"
)

var (
	//nolint:gochecknoglobals
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "url", "handler"},
	)

	//nolint:gochecknoglobals
	HTTPRequestDurationSec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "url", "handler"},
	)
)

//nolint:gochecknoinits
func init() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDurationSec)
}

func handler() http.HandlerFunc {
	promHandler := promhttp.Handler()

	return func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.NotFound(writer, req)

			return
		}

		promHandler.ServeHTTP(writer, req)
	}
}

func StartServer(port int) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           handler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).
				Int("port", port).
				Msg("failed to start metrics server")
		}
	}()

	log.Debug().Int("port", port).Msg("metrics server started")

	return server
}
