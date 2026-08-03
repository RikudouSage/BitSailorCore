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

	FetchAuthRequest(ctx context.Context, session *result.Session, id uuid.UUID) (*result.AuthRequest, error)
	RespondToAuthRequest(ctx context.Context, session *result.Session, request *result.AuthRequest, approved bool) error
}

type auth struct {
	identityURL *url.URL
	apiURL      *url.URL
	httpClient  *http.Client
	deviceID    uuid.UUID

	debugLogs bool

	now func() time.Time
}

func newAuth(
	identityURL *url.URL,
	apiURL *url.URL,
	httpClient *http.Client,
	deviceID uuid.UUID,
	debugLogs bool,
) *auth {
	return &auth{
		identityURL: identityURL,
		apiURL:      apiURL,
		httpClient:  httpClient,
		deviceID:    deviceID,
		now:         time.Now,
		debugLogs:   debugLogs,
	}
}

func (receiver *auth) getTokenURL() *url.URL {
	return urlWithPath(receiver.identityURL, "/identity/connect/token")
}
