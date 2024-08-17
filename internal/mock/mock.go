package mock

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func StartServer(port string, handler http.HandlerFunc) *http.Server {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).
				Str("port", port).
				Msg("failed to start mock server")
		}
	}()

	log.Debug().Str("port", port).Msg("mock server started")

	return server
}
