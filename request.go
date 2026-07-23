package bitwarden

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.chrastecky.dev/bitwarden-client/bitwarden/internal"
	internalHttp "go.chrastecky.dev/bitwarden-client/bitwarden/internal/http"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func request[TResponse any](
	ctx context.Context,
	httpClient *http.Client,
	method string,
	url *url.URL,
	body any,
	session *result.Session,
) (TResponse, error) {
	var requestBody io.Reader
	var out TResponse

	if body != nil {
		var processedBody []byte
		if strBody, ok := body.(string); ok {
			processedBody = []byte(strBody)
		} else if bytesBody, ok := body.([]byte); ok {
			processedBody = bytesBody
		} else {
			var err error
			processedBody, err = json.Marshal(body)
			if err != nil {
				return out, fmt.Errorf("failed marshalling body to json: %w", err)
			}
		}
		requestBody = bytes.NewReader(processedBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), requestBody)
	if err != nil {
		return out, fmt.Errorf("failed creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Bitwarden-Client-Version", internal.BitwardenVersion)
	req.Header.Set("Bitwarden-Client-Name", internal.DeviceName)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if session != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", session.Auth.TokenType, session.Auth.AccessToken))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("failed sending request: %w", err)
	}
	defer internalHttp.DrainResponse(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ = io.ReadAll(resp.Body)
		return out, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		if errors.Is(err, io.EOF) {
			return out, nil
		}
		return out, fmt.Errorf("failed decoding response: %w", err)
	}

	return out, nil
}
