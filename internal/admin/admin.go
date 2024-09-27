package admin

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	path = "/admin"

	readHeaderTimeout = 5 * time.Second
)

func handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.NotFound(writer, req)

			return
		}

		writer.WriteHeader(http.StatusOK)

		_, err := writer.Write([]byte("OK"))
		if err != nil {
			log.Warn().Err(err).Msg("failed to write response")
		}
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
