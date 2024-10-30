package mock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lameaux/mox/internal/mock/predefined"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func NewHandler(configPath string, accessLog bool) (http.Handler, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))

	var mappings []*Mapping

	if configPath != "" {
		mappings, err = loadMappings(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load mappings from %v: %w", configPath, err)
		}
		log.Debug().Msg("mappings loaded successfully")
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		renderMapping(w, r, mappings, accessLog)
	}

	instrumented := otelhttp.NewHandler(
		http.HandlerFunc(f),
		"mox",
		otelhttp.WithMeterProvider(provider),
	)

	return instrumented, nil
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
		found.render(req.Context(), writer)
	} else {
		handlerName = predefined.Render(writer, req)
	}

	latency := time.Since(startTime)

	addMetricLabels(req, handlerName)

	if accessLog {
		log.Debug().
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Str("handler", handlerName).
			Dur("latency", latency).
			Msg("access log")
	}
}

func addMetricLabels(r *http.Request, handlerName string) {
	labeler, _ := otelhttp.LabelerFromContext(r.Context())
	labeler.Add(
		attribute.KeyValue{
			Key:   "url",
			Value: attribute.StringValue(r.URL.Path),
		},
		attribute.KeyValue{
			Key:   "handler",
			Value: attribute.StringValue(handlerName),
		},
	)
}
