package param

import (
	"net/http"
	"strconv"
)

func IntQueryParam(r *http.Request, param string, defaultValue int) (int, error) {
	queryParams := r.URL.Query()

	val := queryParams.Get(param)
	if val == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(val) //nolint:wrapcheck
}
