package bitwarden

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

type Auth interface {
	LoginPassword(ctx context.Context, email, password string, twoFaCode *string) (*result.Session, error)
	LoginApiKey(ctx context.Context, clientID, clientSecret string) (*result.Session, error)
	RefreshToken(ctx context.Context, session *result.Session) error
	UnlockSession(ctx context.Context, session *result.Session, email, password string) error
}

type auth struct {
	identityURL *url.URL
	httpClient  *http.Client
	deviceID    uuid.UUID

	now func() time.Time
}

func newAuth(
	identityURL *url.URL,
	httpClient *http.Client,
	deviceID uuid.UUID,
) *auth {
	return &auth{
		identityURL: identityURL,
		httpClient:  httpClient,
		deviceID:    deviceID,
		now:         time.Now,
	}
}

func (receiver *auth) getTokenURL() *url.URL {
	return urlWithPath(receiver.identityURL, "/identity/connect/token")
}
