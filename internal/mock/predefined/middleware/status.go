package middleware

import (
	"fmt"
	"net/http"

	"github.com/lameaux/mox/internal/mock/predefined/param"
)

const (
	queryNameCode = "code"
)

func Status(writer http.ResponseWriter, reader *http.Request) error {
	statusCode, err := param.IntQueryParam(reader, queryNameCode, defaultZero)
	if err != nil {
		return fmt.Errorf("invalid value for %s: %w", queryNameCode, err)
	}

	if statusCode != defaultZero {
		writer.WriteHeader(statusCode)
	}

	return nil
}
