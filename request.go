package bitwarden

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	internalHttp "go.chrastecky.dev/bitwarden-client/bitwarden/internal/http"
)

func request[TResponse any](ctx context.Context, httpClient *http.Client, method string, url *url.URL, body any) (TResponse, error) {
	var processedBody []byte
	var result TResponse

	if strBody, ok := body.(string); ok {
		processedBody = []byte(strBody)
	} else if bytesBody, ok := body.([]byte); ok {
		processedBody = bytesBody
	} else {
		var err error
		processedBody, err = json.Marshal(body)
		if err != nil {
			return result, fmt.Errorf("failed marshalling body to json: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), bytes.NewReader(processedBody))
	if err != nil {
		return result, fmt.Errorf("failed creating request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed sending request: %w", err)
	}
	defer internalHttp.DrainResponse(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ = io.ReadAll(resp.Body)
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed decoding response: %w", err)
	}

	return result, nil
}
