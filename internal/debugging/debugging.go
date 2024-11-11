package debugging

import (
	"errors"
	_ "expvar" // exposes /debug/vars
	"fmt"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
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

	log.Debug().Int("port", port).Msg("debugging server started")

	return server
}
