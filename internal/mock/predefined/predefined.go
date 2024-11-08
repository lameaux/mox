package predefined

import (
	"net/http"

	"github.com/lameaux/mox/internal/mock/predefined/handler"
	"github.com/lameaux/mox/internal/mock/predefined/middleware"
)

//nolint:gochecknoglobals
var mapping = map[string]http.HandlerFunc{
	"/mox/echo":       handler.Echo,
	"/mox/proxy/http": handler.HTTPProxy,
	"/mox/uuid":       handler.UUIDString,
}

type middlewareFunc func(http.ResponseWriter, *http.Request) error

//nolint:gochecknoglobals
var middlewares = []middlewareFunc{
	middleware.Sleep,
	middleware.Status,
}

const (
	responseNotFound = "not_found"
	middlewareError  = "middleware_error"
)

func Render(writer http.ResponseWriter, req *http.Request) string {
	for _, m := range middlewares {
		if err := m(writer, req); err != nil {
			handler.SendError(writer, http.StatusInternalServerError, err)

			return middlewareError
		}
	}

	handlerFunc, ok := mapping[req.URL.Path]
	if !ok {
		http.NotFound(writer, req)

		return responseNotFound
	}

	handlerFunc(writer, req)

	return req.URL.Path
}
