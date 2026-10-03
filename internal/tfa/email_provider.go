package tfa

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/dto"
	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
)

type EmailProvider interface {
	Provider
}

type emailProvider struct {
	httpClient *http.Client
	apiURL     *url.URL
	deviceID   uuid.UUID
	debug      bool
}

type emailTFARequest struct {
	Email              string    `json:"email"`
	MasterPasswordHash string    `json:"masterPasswordHash"`
	DeviceID           uuid.UUID `json:"deviceIdentifier"`
}

func NewEmailProvider(
	httpClient *http.Client,
	apiURL *url.URL,
	deviceID uuid.UUID,
	debug bool,
) EmailProvider {
	return &emailProvider{
		httpClient: httpClient,
		apiURL:     apiURL,
		deviceID:   deviceID,
		debug:      debug,
	}
}

func (receiver *emailProvider) Kind() dto.TFAKind {
	return dto.TFAKindEmail
}

func (receiver *emailProvider) Name() string {
	return "email"
}

func (receiver *emailProvider) Initialize(ctx context.Context, email, hash string) error {
	body := &emailTFARequest{
		Email:              email,
		MasterPasswordHash: hash,
		DeviceID:           receiver.deviceID,
	}

	_, err := internalHttp.DoRequest[any](
		ctx,
		receiver.httpClient,
		http.MethodPost,
		internalHttp.UrlWithPath(receiver.apiURL, "/two-factor/send-email-login"),
		body,
		nil,
		receiver.debug,
	)

	if err != nil {
		return fmt.Errorf("failed issuing an email tfa request: %w", err)
	}

	return nil
}

func (receiver *emailProvider) ModifyRequest(request Request) error {
	request.SetProvider(new(receiver.Kind()))
	request.SetRemember(new(false))

	return nil
}
