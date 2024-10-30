package predefined

import (
	"io"
	"net/http"
	"time"

	"github.com/lameaux/mox/internal/config"
	"github.com/lameaux/mox/internal/httpclient"
	"github.com/rs/zerolog/log"
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

	req, err := http.NewRequestWithContext(httpReq.Context(), http.MethodGet, url, nil)
	if err != nil {
		log.Warn().Err(err).Msg("failed to generate request")
		sendError(httpWriter, http.StatusBadGateway, err)

		return
	}

	conf := config.HTTPClient{
		Timeout: time.Duration(timeout) * time.Second,
	}
	client := httpclient.New(conf)

	resp, err := client.Do(req)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get content")
		sendError(httpWriter, http.StatusBadGateway, err)

		return
	}

	defer resp.Body.Close()

	httpWriter.WriteHeader(resp.StatusCode)

	for k, v := range resp.Header {
		httpWriter.Header().Set(k, v[0])
	}

	_, err = io.Copy(httpWriter, resp.Body)
	if err != nil {
		log.Warn().Err(err).Msg("failed to write response")
	}
}
