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
		if baseURL == "" {
			bwClient.identityURL = nil
			bwClient.apiURL = nil
			bwClient.sendURL = nil
			return nil
		}

		parsed, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("failed parsing base url: %w", err)
		}
		bwClient.identityURL = parsed
		bwClient.apiURL = parsed
		bwClient.sendURL = parsed
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

func WithSendURL(sendURL string) Option {
	return func(bwClient *client) error {
		if sendURL == "" {
			bwClient.sendURL = nil
			return nil
		}

		parsed, err := url.Parse(sendURL)
		if err != nil {
			return fmt.Errorf("failed parsing send url: %w", err)
		}
		bwClient.sendURL = parsed
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
