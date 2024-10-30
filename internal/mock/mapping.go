package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lameaux/mox/internal/config"
	"github.com/lameaux/mox/internal/httpclient"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	headerContentType = "Content-Type"
	mappingsDir       = "mappings"
	filesDir          = "files"
	templatesDir      = "templates"
)

type Mapping struct {
	ConfigPath string `json:"-"`
	FileName   string `json:"-"`

	Request struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	} `json:"request"`

	Response *struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`

		Body         string         `json:"body"`
		BodyJSON     map[string]any `json:"json"`
		BodyFile     string         `json:"file"`
		RenderedBody []byte         `json:"-"`

		Template         string             `json:"template"`
		TemplateFile     string             `json:"templateFile"`
		RenderedTemplate *template.Template `json:"-"`
	} `json:"response"`

	Proxy *struct {
		Method  string `json:"method"`
		URL     string `json:"url"`
		Timeout int    `json:"timeout"`
	} `json:"proxy"`
}

func (m *Mapping) filePath() string {
	return m.ConfigPath + "/" + m.FileName
}

func loadMappings(configPath string) ([]*Mapping, error) {
	var mappings []*Mapping

	mappingsPath := configPath + "/" + mappingsDir

	err := filepath.Walk(mappingsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to load file %v: %w", path, err)
		}

		if info.IsDir() {
			// skip dirs
			return nil
		}

		statInfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat file %v: %w", path, err)
		}

		if statInfo.IsDir() {
			// skip symlinks that points to dir
			return nil
		}

		log.Debug().Str("path", path).Msg("found mapping file")

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %v: %w", path, err)
		}
		defer file.Close()

		var newMapping Mapping

		d := json.NewDecoder(file)

		if err = d.Decode(&newMapping); err != nil {
			return fmt.Errorf("failed to open file %v: %w", path, err)
		}

		newMapping.ConfigPath = configPath
		newMapping.FileName = filepath.Base(path)

		err = newMapping.prerender()
		if err != nil {
			return fmt.Errorf("failed to prerender file %v: %w", path, err)
		}

		mappings = append(mappings, &newMapping)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load mappings: %w", err)
	}

	return mappings, nil
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
	if m.Response == nil {
		return nil
	}

	if m.Response.Headers == nil {
		m.Response.Headers = make(map[string]string)
	}

	if m.Response.BodyJSON != nil {
		return m.prerenderJSON()
	}

	if m.Response.BodyFile != "" {
		return m.prerenderFile()
	}

	if m.Response.Template != "" {
		return m.prerenderTemplate()
	}

	if m.Response.TemplateFile != "" {
		return m.prerenderTemplateFile()
	}

	// read from body by default
	m.Response.RenderedBody = []byte(m.Response.Body)

	return nil
}

func (m *Mapping) prerenderJSON() error {
	body, err := json.Marshal(m.Response.BodyJSON)
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

func (m *Mapping) prerenderFile() error {
	path := m.ConfigPath + "/" + filesDir + "/" + m.Response.BodyFile

	body, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %v: %w", path, err)
	}

	m.Response.RenderedBody = body

	return nil
}

func (m *Mapping) prerenderTemplate() error {
	t, err := template.New(m.filePath()).Funcs(funcMap()).Parse(m.Response.Template)
	if err != nil {
		return fmt.Errorf("failed to parse template %v: %w", m.filePath(), err)
	}

	m.Response.RenderedTemplate = t

	return nil
}

func (m *Mapping) prerenderTemplateFile() error {
	path := m.ConfigPath + "/" + templatesDir + "/" + m.Response.TemplateFile

	t, err := template.New(filepath.Base(path)).Funcs(funcMap()).ParseFiles(path)
	if err != nil {
		return fmt.Errorf("failed to parse template %v: %w", path, err)
	}

	m.Response.RenderedTemplate = t

	return nil
}
func (m *Mapping) render(ctx context.Context, writer http.ResponseWriter) {
	var err error

	if m.Proxy != nil {
		err = m.renderProxy(ctx, writer)
	} else {
		err = m.renderResponse(writer)
	}

	if err != nil {
		log.Warn().Err(err).Str("mapping", m.filePath()).Msg("failed to write response")
	}
}

func (m *Mapping) renderProxy(ctx context.Context, writer http.ResponseWriter) error {
	client := httpclient.New(
		config.HTTPClient{
			Timeout: time.Duration(m.Proxy.Timeout) * time.Second,
		},
	)

	return httpclient.Proxy(ctx, m.Proxy.Method, m.Proxy.URL, client, writer)
}

func (m *Mapping) renderResponse(writer http.ResponseWriter) error {
	for k, v := range m.Response.Headers {
		writer.Header().Set(k, v)
	}

	writer.WriteHeader(m.Response.Status)

	var err error
	if m.Response.RenderedTemplate != nil {
		err = m.Response.RenderedTemplate.Execute(writer, nil)
	} else {
		_, err = writer.Write(m.Response.RenderedBody)
	}

	return err
}
