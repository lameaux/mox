package handler

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func SendError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)

	_, err2 := w.Write([]byte(err.Error()))
	if err2 != nil {
		log.Warn().Err(err2).Msg("failed to write response")
	}
}
