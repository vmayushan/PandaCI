package retryClient

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type RetryRoundTripper struct {
	Base       http.RoundTripper
	MaxRetries int
	Headers    map[string]string
}

func (hrt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	for k, v := range hrt.Headers {
		req.Header.Add(k, v)
	}

	retries := hrt.MaxRetries
	// Set default values
	if hrt.MaxRetries == 0 {
		retries = 5
	}

	for i := 1; i <= retries; i++ {
		resp, err := hrt.Base.RoundTrip(req)
		if err == nil && resp.StatusCode < 400 {
			return resp, err
		}

		if err != nil {
			// log error
			log.Error().Err(err).Msg("failed to make request")
		}

		// exponential backoff
		time.Sleep(time.Duration(i*i) * time.Second)
	}

	return hrt.Base.RoundTrip(req)
}
