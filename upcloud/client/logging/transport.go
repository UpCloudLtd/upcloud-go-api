package logging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type LoggingTransport struct {
	Logger    *Logger
	Transport http.RoundTripper
}

func (t LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	if t.Logger != nil {
		var body []byte

		if req.GetBody != nil {
			bodyReader, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("failed to get request body: %w", err)
			}
			body, err = io.ReadAll(bodyReader)
			if err != nil {
				return nil, fmt.Errorf("failed to read request body: %w", err)
			}
		}
		t.Logger.LogRequest(req, body)
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if t.Logger != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, fmt.Errorf("failed to read response body: %w", err)
		}
		resp.Body = io.NopCloser(bytes.NewReader(body))

		t.Logger.LogResponse(resp, body)
	}

	return resp, err
}
