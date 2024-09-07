package mock

import (
	"net/http"
	"strconv"
	"time"
)

var mapping = map[string]http.HandlerFunc{
	"/mox/sleep": moxSleep,
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
	queryParams := r.URL.Query()
	seconds := queryParams.Get("seconds")
	if seconds == "" {
		seconds = "0"
	}

	sleepSec, err := strconv.Atoi(seconds)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}

	time.Sleep(time.Duration(sleepSec) * time.Second)
}

func sendError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
