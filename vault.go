package bitwarden

import (
	"context"
	"net/http"
	"net/url"

	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

type Vault interface {
	Sync(ctx context.Context, session *result.Session) (*result.Sync, error)
}

type vault struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func newVault(baseURL *url.URL, httpClient *http.Client) *vault {
	return &vault{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}
