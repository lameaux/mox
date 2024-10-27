package pprof

import (
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof" //nolint:gosec // Import for pprof, only enabled via CLI flag
	"time"

	"github.com/rs/zerolog/log"
)

const (
	readHeaderTimeout = 5 * time.Second
)

func StartServer(port int) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).
				Int("port", port).
				Msg("failed to start pprof server")
		}
	}()

	log.Debug().Int("port", port).Msg("pprof server started")

	return server
}
