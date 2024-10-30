package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/lameaux/mox/internal/config"
	"github.com/rs/zerolog/log"
)

const defaultMaxIdleConnsPerHost = 100

func New(conf config.HTTPClient) *http.Client {
	maxIdleConnsPerHost := defaultMaxIdleConnsPerHost
	if conf.MaxIdleConnsPerHost > 0 {
		maxIdleConnsPerHost = conf.MaxIdleConnsPerHost
	}

	log.Debug().
		Bool("disableKeepAlive", conf.DisableKeepAlive).
		Dur("timeout", conf.Timeout).
		Int("maxIdleConnsPerHost", maxIdleConnsPerHost).
		Msg("creating http client")

	transport := &http.Transport{
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		DisableKeepAlives:   conf.DisableKeepAlive,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   conf.Timeout,
	}

	if conf.DisableFollowRedirects {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

func Proxy(ctx context.Context, method string, url string, client *http.Client, httpWriter http.ResponseWriter) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("failed to generate request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get content: %w", err)
	}

	defer resp.Body.Close()

	httpWriter.WriteHeader(resp.StatusCode)

	for k, v := range resp.Header {
		httpWriter.Header().Set(k, v[0])
	}

	_, err = io.Copy(httpWriter, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}
