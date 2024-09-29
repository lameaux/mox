package predefined

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

//nolint:gochecknoglobals
var mapping = map[string]http.HandlerFunc{
	"/mox/sleep":        sleep,
	"/mox/sleep/random": sleepRandom,
	"/mox/rabbitmq":     rabbitmq,
}

const (
	responseNotFound = "not_found"
)

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

	seconds := queryParams.Get(param)
	if seconds == "" {
		seconds = "0"
	}

	return strconv.Atoi(seconds) //nolint:wrapcheck
}

func sendError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)

	_, err2 := w.Write([]byte(err.Error()))
	if err2 != nil {
		log.Warn().Err(err2).Msg("failed to write response")
	}
}
