package predefined

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	queryNameMs = "ms"
)

func sleep(w http.ResponseWriter, r *http.Request) {
	sleepMillis, err := getIntQueryParam(r, queryNameMs)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}

	time.Sleep(time.Duration(sleepMillis) * time.Millisecond)
}

func sleepRandom(w http.ResponseWriter, r *http.Request) {
	sleepMillis, err := getIntQueryParam(r, queryNameMs)
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
	}

	randomSleep := rand.Intn(sleepMillis + 1) //nolint:gosec

	time.Sleep(time.Duration(randomSleep) * time.Millisecond)
}
