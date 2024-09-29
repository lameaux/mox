package predefined

import (
	"errors"
	"net/http"
)

var errInvalidRequest = errors.New("invalid request")

func rabbitmq(w http.ResponseWriter, _ *http.Request) {
	sendError(w, http.StatusBadRequest, errInvalidRequest)
}
