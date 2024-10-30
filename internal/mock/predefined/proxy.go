package predefined

import (
	"net/http"
	"time"

	"github.com/lameaux/mox/internal/config"
	"github.com/lameaux/mox/internal/httpclient"
)

func httpProxy(httpWriter http.ResponseWriter, httpReq *http.Request) {
	url := httpReq.URL.Query().Get("url")
	if url == "" {
		sendError(httpWriter, http.StatusBadRequest, errInvalidRequest)

		return
	}

	timeout, err := getIntQueryParam(httpReq, "timeout")
	if err != nil {
		sendError(httpWriter, http.StatusBadRequest, errInvalidRequest)

		return
	}

	conf := config.HTTPClient{
		Timeout: time.Duration(timeout) * time.Second,
	}
	client := httpclient.New(conf)

	err = httpclient.Proxy(httpReq.Context(), http.MethodGet, url, client, httpWriter)
	if err != nil {
		sendError(httpWriter, http.StatusBadGateway, err)

		return
	}
}
