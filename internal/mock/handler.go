package mock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lameaux/mox/internal/metrics"
	"github.com/rs/zerolog/log"
)

func NewHandler(configPath string, accessLog bool) (http.HandlerFunc, error) {
	mappings, err := loadMappings(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load mappings from %v: %w", configPath, err)
	}

	log.Debug().Msg("mappings loaded successfully")

	f := func(w http.ResponseWriter, r *http.Request) {
		renderMapping(w, r, mappings, accessLog)
	}

	return f, nil
}

func renderMapping(writer http.ResponseWriter, req *http.Request, mappings []*Mapping, accessLog bool) {
	startTime := time.Now()

	var found *Mapping

	for _, m := range mappings {
		if m.matches(req) {
			found = m

			break
		}
	}

	var handlerName string

	if found != nil {
		handlerName = found.filePath()
		found.render(writer)
	} else {
		handlerName = renderPredefined(writer, req)
	}

	latency := time.Since(startTime)

	metrics.HTTPRequestsTotal.WithLabelValues(
		req.Method,
		req.URL.String(),
		handlerName,
	).Inc()

	metrics.HTTPRequestDurationSec.WithLabelValues(
		req.Method,
		req.URL.String(),
		handlerName,
	).Observe(latency.Seconds())

	if accessLog {
		log.Debug().
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Str("handler", handlerName).
			Dur("latency", latency).
			Msg("access log")
	}
}
