package predefined

import (
	"net/http"
)

func rabbitmq(w http.ResponseWriter, _ *http.Request) {
	sendError(w, http.StatusBadRequest, errInvalidRequest)
}
