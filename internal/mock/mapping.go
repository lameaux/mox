package mock

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	headerContentType = "Content-Type"
	mappingsDir       = "mappings"
	filesDir          = "files"
	templatesDir      = "templates"
)

type Mapping struct {
	ConfigPath string
	FileName   string

	Request struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	}

	Response struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`

		Body         string         `json:"body"`
		BodyJson     map[string]any `json:"json"`
		BodyFile     string         `json:"file"`
		RenderedBody []byte

		Template         string `json:"template"`
		TemplateFile     string `json:"templateFile"`
		RenderedTemplate *template.Template
	}
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
			// skip
			return nil
		}

		log.Debug().Str("path", path).Msg("found mapping file")

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %v: %w", path, err)
		}
		defer file.Close()

		var m Mapping

		d := json.NewDecoder(file)
		err = d.Decode(&m)
		if err != nil {
			return fmt.Errorf("failed to open file %v: %w", path, err)
		}

		m.ConfigPath = configPath
		m.FileName = filepath.Base(path)

		err = m.prerender()
		if err != nil {
			return fmt.Errorf("failed to prerender file %v: %w", path, err)
		}

		mappings = append(mappings, &m)

		return nil
	})

	if err != nil {
		return nil, err
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
	if m.Response.Headers == nil {
		m.Response.Headers = make(map[string]string)
	}

	if m.Response.BodyJson != nil {
		return m.prerenderJson()
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

func (m *Mapping) prerenderJson() error {
	body, err := json.Marshal(m.Response.BodyJson)
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

func (m *Mapping) render(w http.ResponseWriter) {
	for k, v := range m.Response.Headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(m.Response.Status)

	var err error
	if m.Response.RenderedTemplate != nil {
		err = m.Response.RenderedTemplate.Execute(w, nil)
	} else {
		_, err = w.Write(m.Response.RenderedBody)
	}

	if err != nil {
		log.Warn().Err(err).Str("mapping", m.filePath()).Msg("failed to write response")
	}
}
