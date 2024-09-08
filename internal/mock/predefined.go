package mock

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var mapping = map[string]http.HandlerFunc{
	"/mox/sleep":        moxSleep,
	"/mox/sleep/random": moxSleepRandom,
}

func renderPredefined(w http.ResponseWriter, r *http.Request) string {
	fn, ok := mapping[r.URL.Path]

	if ok {
		fn(w, r)
		return r.URL.Path
	}

	http.NotFound(w, r)
	return "not_found"
}

func moxSleep(w http.ResponseWriter, r *http.Request) {
	sleepMillis, err := getIntQueryParam(r, "ms")
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}

	time.Sleep(time.Duration(sleepMillis) * time.Millisecond)
}

func moxSleepRandom(w http.ResponseWriter, r *http.Request) {
	sleepMillis, err := getIntQueryParam(r, "ms")
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}

	randomSleep := rand.Intn(sleepMillis + 1)

	time.Sleep(time.Duration(randomSleep) * time.Millisecond)
}

func getIntQueryParam(r *http.Request, param string) (int, error) {
	queryParams := r.URL.Query()
	seconds := queryParams.Get(param)
	if seconds == "" {
		seconds = "0"
	}

	return strconv.Atoi(seconds)
}

func sendError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
