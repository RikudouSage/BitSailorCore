package bitwarden

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

type Auth interface {
	LoginPassword(ctx context.Context, email, password string) (*result.Session, error)
	LoginApiKey(ctx context.Context, clientID, clientSecret string) (*result.Session, error)
	UnlockSession(ctx context.Context, session *result.Session, email, password string) error
}

type auth struct {
	baseURL    *url.URL
	httpClient *http.Client
	deviceID   uuid.UUID

	now func() time.Time
}

func newAuth(
	baseURL *url.URL,
	httpClient *http.Client,
	deviceID uuid.UUID,
) *auth {
	return &auth{
		baseURL:    baseURL,
		httpClient: httpClient,
		deviceID:   deviceID,
		now:        time.Now,
	}
}

func (receiver *auth) getTokenURL() *url.URL {
	tokenURL := new(*receiver.baseURL)
	tokenURL.Path = "/identity/connect/token"

	return tokenURL
}
