package mock

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	readHeaderTimeout = 5 * time.Second
)

func StartServer(port int, handler http.HandlerFunc) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).
				Int("port", port).
				Msg("failed to start mock server")
		}
	}()

	log.Debug().Int("port", port).Msg("mock server started")

	return server
}
