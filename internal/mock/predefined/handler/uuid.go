package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func UUIDString(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(uuid.NewString()))
	if err != nil {
		log.Warn().Err(err).Msg("failed to write uuid response")
	}
}
