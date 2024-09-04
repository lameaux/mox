package mock

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func StartServer(port int, handler http.HandlerFunc) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).
				Int("port", port).
				Msg("failed to start mock server")
		}
	}()

	log.Debug().Int("port", port).Msg("mock server started")

	return server
}
