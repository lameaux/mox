package mock

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

//nolint:gochecknoglobals
var mapping = map[string]http.HandlerFunc{
	"/mox/sleep":        moxSleep,
	"/mox/sleep/random": moxSleepRandom,
}

const (
	responseNotFound = "not_found"
	queryNameMs      = "ms"
)

func renderPredefined(writer http.ResponseWriter, req *http.Request) string {
	fn, ok := mapping[req.URL.Path]

	if ok {
		fn(writer, req)

		return req.URL.Path
	}

	http.NotFound(writer, req)

	return responseNotFound
}

func moxSleep(w http.ResponseWriter, r *http.Request) {
	sleepMillis, err := getIntQueryParam(r, queryNameMs)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}

	time.Sleep(time.Duration(sleepMillis) * time.Millisecond)
}

func moxSleepRandom(w http.ResponseWriter, r *http.Request) {
	sleepMillis, err := getIntQueryParam(r, queryNameMs)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}

	randomSleep := rand.Intn(sleepMillis + 1) //nolint:gosec

	time.Sleep(time.Duration(randomSleep) * time.Millisecond)
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
