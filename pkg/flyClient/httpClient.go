package flyClient

import (
	"net/http"

	"github.com/pandaci-com/pandaci/pkg/retryClient"
)

type FlyRoundTripper struct {
	Base      http.RoundTripper
	AppName   *string
	MachineID *string
	Headers   map[string]string
}

func (hrt *FlyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	if hrt.AppName != nil {
		req.Header.Add("PandaCI-app-name", *hrt.AppName)
	}

	for key, value := range hrt.Headers {
		req.Header.Add(key, value)
	}

	retryRoundTripper := &retryClient.RetryRoundTripper{
		Base: hrt.Base,
	}

	return retryRoundTripper.RoundTrip(req)
}
