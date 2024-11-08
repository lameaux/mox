package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/lameaux/mox/internal/config"
	"github.com/lameaux/mox/internal/httpclient"
	"github.com/lameaux/mox/internal/mock/predefined/param"
)

const (
	defaultTimeout = 0
)

var errInvalidRequest = errors.New("invalid request")

func HTTPProxy(httpWriter http.ResponseWriter, httpReq *http.Request) {
	url := httpReq.URL.Query().Get("url")
	if url == "" {
		SendError(httpWriter, http.StatusBadRequest, errInvalidRequest)

		return
	}

	timeout, err := param.IntQueryParam(httpReq, "timeout", defaultTimeout)
	if err != nil {
		SendError(httpWriter, http.StatusBadRequest, errInvalidRequest)

		return
	}

	conf := config.HTTPClient{
		Timeout: time.Duration(timeout) * time.Second,
	}
	client := httpclient.New(conf)

	err = httpclient.Proxy(httpReq.Context(), http.MethodGet, url, client, httpWriter)
	if err != nil {
		SendError(httpWriter, http.StatusBadGateway, err)

		return
	}
}
