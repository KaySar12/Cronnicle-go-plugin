package httpHelper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func HTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        10,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func ExecuteRequest(ctx context.Context, method string, url string, body io.Reader, apiKey string, result interface{}, retries ...int) error {
	maxRetries := 1
	if len(retries) > 0 {
		maxRetries = retries[0]
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("X-API-Key", apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	// Retry logic
	for i := 0; i < maxRetries; i++ {
		httpResp, err := HTTPClient().Do(httpReq)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Request timed out, retrying...")
				continue
			}
			return fmt.Errorf("error making request: %v", err)
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("HTTP error: %s", httpResp.Status)
		}

		if result != nil {
			if err := json.NewDecoder(httpResp.Body).Decode(result); err != nil {
				return fmt.Errorf("error decoding response: %v", err)
			}
		}
		return nil
	}
	return fmt.Errorf("failed to get a successful response after retries")
}

func GetHttpRes(url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := HTTPClient().Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		} else {
			if strings.Contains(err.Error(), "no such host") {
				return nil, err
			}
			return nil, err
		}
	}
	if resp.StatusCode == 404 {
		return nil, err
	}

	return resp, nil
}
