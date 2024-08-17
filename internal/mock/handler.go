package mock

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct {
	mappings []Mapping
}

func NewHandler(mappingsPath string) (http.HandlerFunc, error) {
	mappings, err := loadMappings(mappingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load mappings from %v: %w", mappingsPath, err)
	}

	log.Debug().Msg("mappings loaded successfully")

	f := func(w http.ResponseWriter, r *http.Request) {
		for _, m := range mappings {
			if m.matches(r) {
				m.render(w)
				return
			}
		}

		http.NotFound(w, r)
	}

	return f, nil
}

func loadMappings(mappingsPath string) ([]*Mapping, error) {
	var mappings []*Mapping

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

		m.FilePath = path

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
