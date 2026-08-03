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
	Notifications() Notifications

	GeneratePassword(request *PasswordGeneratorRequest) (string, error)
	GeneratePassphrase(request *PassphraseGeneratorRequest) (string, error)
}

type client struct {
	httpClient       *http.Client
	identityURL      *url.URL
	apiURL           *url.URL
	notificationsURL *url.URL
	deviceID         uuid.UUID

	debugLogs bool

	auth          *auth
	vault         *vault
	notifications *notifications
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
		receiver.auth = newAuth(receiver.identityURL, receiver.httpClient, receiver.deviceID, receiver.debugLogs)
	}

	return receiver.auth
}

func (receiver *client) Vault() Vault {
	if receiver.vault == nil {
		receiver.vault = newVault(
			receiver.apiURL,
			receiver.identityURL,
			receiver.httpClient,
			receiver.Auth().(*auth),
		)
	}

	return receiver.vault
}

func (receiver *client) Notifications() Notifications {
	if receiver.notifications == nil {
		receiver.notifications = newNotifications(
			receiver.notificationsURL,
			receiver.httpClient,
			receiver.deviceID,
			receiver.Auth().(*auth),
		)
	}

	return receiver.notifications
}
