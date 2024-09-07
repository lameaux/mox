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
		Render(w, r, mappings, accessLog)
	}

	return f, nil
}

func Render(w http.ResponseWriter, r *http.Request, mappings []*Mapping, accessLog bool) {
	startTime := time.Now()

	var found *Mapping

	for _, m := range mappings {
		if m.matches(r) {
			found = m
			break
		}
	}

	mappingFile := "not found"

	if found != nil {
		mappingFile = found.filePath()
		found.render(w)
	}

	latency := time.Now().Sub(startTime)

	metrics.HttpRequestsTotal.WithLabelValues(
		r.Method,
		r.URL.String(),
		mappingFile,
	).Inc()

	metrics.HttpRequestDurationSec.WithLabelValues(
		r.Method,
		r.URL.String(),
		mappingFile,
	).Observe(latency.Seconds())

	if accessLog {
		log.Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("mapping", mappingFile).
			Dur("latency", latency).
			Msg("access log")
	}

	if found == nil {
		http.NotFound(w, r)
	}
}
