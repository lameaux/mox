package mock

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	headerContentType = "Content-Type"
)

type Mapping struct {
	FilePath string
	Request  struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	}
	Response struct {
		Status       int               `json:"status"`
		Headers      map[string]string `json:"headers"`
		Body         string            `json:"body"`
		JsonBody     map[string]any    `json:"jsonBody"`
		RenderedBody []byte
	}
}

func (m *Mapping) matches(r *http.Request) bool {
	if m.Request.URL != r.URL.Path {
		return false
	}

	if !strings.EqualFold(m.Request.Method, r.Method) {
		return false
	}

	return true
}

func (m *Mapping) prerender() error {
	if m.Response.Headers == nil {
		m.Response.Headers = make(map[string]string)
	}

	if m.Response.JsonBody != nil {
		body, err := json.Marshal(m.Response.JsonBody)
		if err != nil {
			return fmt.Errorf("failed to render jsonBody: %w", err)
		}

		m.Response.RenderedBody = body

		contentTypeSet := false
		for h := range m.Response.Headers {
			if strings.EqualFold(h, headerContentType) {
				contentTypeSet = true
				break
			}
		}

		if !contentTypeSet {
			m.Response.Headers[headerContentType] = "application/json"
		}

		return nil
	}

	m.Response.RenderedBody = []byte(m.Response.Body)

	return nil
}

func (m *Mapping) render(w http.ResponseWriter) {
	for k, v := range m.Response.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(m.Response.Status)

	_, err := w.Write(m.Response.RenderedBody)
	if err != nil {
		log.Warn().Err(err).Str("mapping", m.FilePath).Msg("failed to write response")
	}
}
