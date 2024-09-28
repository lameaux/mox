package admin

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lameaux/mox/internal/banner"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const (
	pathIndex   = "/"
	pathMetrics = "/metrics"
	pathAPI     = "/api"

	readHeaderTimeout = 5 * time.Second
)

func handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		h := handlerByRequest(req)
		h.ServeHTTP(writer, req)
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
				Msg("failed to start admin server")
		}
	}()

	log.Debug().Int("port", port).Msg("admin server started")

	return server
}

func handlerByRequest(req *http.Request) http.Handler {
	switch req.URL.Path {
	case pathIndex:
		return IndexHandler()
	case pathAPI:
		return APIHandler()
	case pathMetrics:
		return promhttp.Handler()
	}

	return http.NotFoundHandler()
}

func IndexHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)

		_, err := writer.Write([]byte(banner.Banner))
		if err != nil {
			log.Warn().Err(err).Msg("failed to write response")
		}
	})
}

func APIHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)

		_, err := writer.Write([]byte("API"))
		if err != nil {
			log.Warn().Err(err).Msg("failed to write response")
		}
	})
}
