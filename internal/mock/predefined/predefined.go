package predefined

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

//nolint:gochecknoglobals
var mapping = map[string]http.HandlerFunc{
	"/mox/sleep":        sleep,
	"/mox/sleep/random": sleepRandom,
	"/mox/rabbitmq":     rabbitmq,
	"/mox/proxy/http":   httpProxy,
}

const (
	responseNotFound = "not_found"
)

var errInvalidRequest = errors.New("invalid request")

func Render(writer http.ResponseWriter, req *http.Request) string {
	handler, ok := mapping[req.URL.Path]
	if !ok {
		http.NotFound(writer, req)

		return responseNotFound
	}

	handler(writer, req)

	return req.URL.Path
}

func getIntQueryParam(r *http.Request, param string) (int, error) {
	queryParams := r.URL.Query()

	val := queryParams.Get(param)
	if val == "" {
		val = "0"
	}

	return strconv.Atoi(val) //nolint:wrapcheck
}

func sendError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)

	_, err2 := w.Write([]byte(err.Error()))
	if err2 != nil {
		log.Warn().Err(err2).Msg("failed to write response")
	}
}
