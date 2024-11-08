package handler

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Echo(httpWriter http.ResponseWriter, httpReq *http.Request) {
	body := httpReq.URL.Query().Get("body")
	if body == "" {
		SendError(httpWriter, http.StatusBadRequest, errInvalidRequest)

		return
	}

	_, err := httpWriter.Write([]byte(body))
	if err != nil {
		log.Warn().Err(err).Msg("failed to write echo response")
	}
}
