package bitwarden

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type Client interface {
	Auth() Auth
	Vault() Vault

	GeneratePassword(request *PasswordGeneratorRequest) (string, error)
}

type client struct {
	httpClient  *http.Client
	identityURL *url.URL
	apiURL      *url.URL
	deviceID    uuid.UUID

	auth  *auth
	vault *vault
}

func NewClient(options ...Option) (Client, error) {
	bwClient := &client{}

	for _, option := range options {
		if err := option(bwClient); err != nil {
			return nil, fmt.Errorf("failed applying an option: %w", err)
		}
	}

	if err := bwClient.provideDefaultsAndValidate(); err != nil {
		return nil, fmt.Errorf("failed validating options: %w", err)
	}

	return bwClient, nil
}

func (receiver *client) Auth() Auth {
	if receiver.auth == nil {
		receiver.auth = newAuth(receiver.identityURL, receiver.httpClient, receiver.deviceID)
	}

	return receiver.auth
}

func (receiver *client) Vault() Vault {
	if receiver.vault == nil {
		receiver.vault = newVault(receiver.apiURL, receiver.httpClient)
	}

	return receiver.vault
}
