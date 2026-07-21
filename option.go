package bitwarden

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type Option func(bwClient *client) error

func WithBaseURL(baseURL string) Option {
	return func(bwClient *client) error {
		parsed, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("failed parsing base url: %w", err)
		}
		bwClient.identityURL = parsed
		return nil
	}
}

func WithAPIURL(apiURL string) Option {
	return func(bwClient *client) error {
		parsed, err := url.Parse(apiURL)
		if err != nil {
			return fmt.Errorf("failed parsing api url: %w", err)
		}
		bwClient.apiURL = parsed
		return nil
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(bwClient *client) error {
		bwClient.httpClient = httpClient
		return nil
	}
}

func WithDeviceID(deviceID uuid.UUID) Option {
	return func(bwClient *client) error {
		bwClient.deviceID = deviceID
		return nil
	}
}
