package bitwarden

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type urlConfig struct {
	baseURL     *url.URL
	identityURL *url.URL
	apiURL      *url.URL
}

var specialURLs = []urlConfig{
	{
		baseURL:     lo.Must(url.Parse("https://bitwarden.com")),
		identityURL: lo.Must(url.Parse("https://vault.bitwarden.com")),
		apiURL:      lo.Must(url.Parse("https://api.bitwarden.com")),
	},
	{
		baseURL:     lo.Must(url.Parse("https://api.bitwarden.com")),
		identityURL: lo.Must(url.Parse("https://vault.bitwarden.eu")),
		apiURL:      lo.Must(url.Parse("https://api.bitwarden.eu")),
	},
}

func normalizeBaseURL(baseURL *url.URL) *urlConfig {
	for _, specialURL := range specialURLs {
		if specialURL.baseURL.String() == baseURL.String() {
			return &specialURL
		}
	}

	return &urlConfig{
		baseURL:     baseURL,
		identityURL: baseURL,
		apiURL:      urlWithPath(baseURL, "/api"),
	}
}

type Option func(bwClient *client) error

func WithBaseURL(baseURL string) Option {
	return func(bwClient *client) error {
		if baseURL == "" {
			bwClient.identityURL = nil
			bwClient.apiURL = nil
			return nil
		}

		parsed, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("failed parsing base url: %w", err)
		}
		normalized := normalizeBaseURL(parsed)

		bwClient.identityURL = normalized.identityURL
		bwClient.apiURL = normalized.apiURL
		return nil
	}
}

func WithIdentityURL(identityURL string) Option {
	return func(bwClient *client) error {
		if identityURL == "" {
			bwClient.identityURL = nil
			return nil
		}

		parsed, err := url.Parse(identityURL)
		if err != nil {
			return fmt.Errorf("failed parsing identity url: %w", err)
		}
		bwClient.identityURL = parsed
		return nil
	}
}

func WithAPIURL(apiURL string) Option {
	return func(bwClient *client) error {
		if apiURL == "" {
			bwClient.apiURL = nil
			return nil
		}

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
