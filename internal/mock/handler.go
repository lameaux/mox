package mock

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func NewHandler(configPath string, accessLog bool) (http.HandlerFunc, error) {
	mappings, err := loadMappings(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load mappings from %v: %w", configPath, err)
	}

	log.Debug().Msg("mappings loaded successfully")

	f := func(w http.ResponseWriter, r *http.Request) {
		for _, m := range mappings {
			if m.matches(r) {
				m.render(w)

				if accessLog {
					log.Debug().
						Str("method", r.Method).
						Str("url", r.URL.String()).
						Str("mapping", m.filePath()).
						Msg("matched")
				}
				return
			}
		}

		if accessLog {
			log.Debug().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Msg("not matched")
		}

		http.NotFound(w, r)
	}

	return f, nil
}
