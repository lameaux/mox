package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/lameaux/mox/internal/mock/predefined/param"
)

const (
	queryNameSleep       = "sleep"
	queryNameSleepRandom = "sleepRandom"
)

func Sleep(_ http.ResponseWriter, reader *http.Request) error {
	sleepMillis, err := param.IntQueryParam(reader, queryNameSleep, defaultZero)
	if err != nil {
		return fmt.Errorf("invalid value for %s: %w", queryNameSleep, err)
	}

	if sleepMillis != defaultZero {
		time.Sleep(time.Duration(sleepMillis) * time.Millisecond)

		return nil
	}

	sleepRandomMillis, err := param.IntQueryParam(reader, queryNameSleepRandom, defaultZero)
	if err != nil {
		return fmt.Errorf("invalid value for %s: %w", queryNameSleepRandom, err)
	}

	if sleepRandomMillis != defaultZero {
		randomSleep := rand.Intn(sleepMillis + 1) //nolint:gosec

		time.Sleep(time.Duration(randomSleep) * time.Millisecond)

		return nil
	}

	return nil
}
