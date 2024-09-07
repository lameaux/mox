package mock

import (
	"fmt"
	"github.com/lameaux/mox/internal/metrics"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
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

func renderMapping(w http.ResponseWriter, r *http.Request, mappings []*Mapping, accessLog bool) {
	startTime := time.Now()

	var found *Mapping

	for _, m := range mappings {
		if m.matches(r) {
			found = m
			break
		}
	}

	var handlerName string

	if found != nil {
		handlerName = found.filePath()
		found.render(w)
	} else {
		handlerName = renderPredefined(w, r)
	}

	latency := time.Now().Sub(startTime)

	metrics.HttpRequestsTotal.WithLabelValues(
		r.Method,
		r.URL.String(),
		handlerName,
	).Inc()

	metrics.HttpRequestDurationSec.WithLabelValues(
		r.Method,
		r.URL.String(),
		handlerName,
	).Observe(latency.Seconds())

	if accessLog {
		log.Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("handler", handlerName).
			Dur("latency", latency).
			Msg("access log")
	}
}
