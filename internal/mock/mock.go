package mock

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/static" {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

func StartServer(port string) *http.Server {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler(),
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
