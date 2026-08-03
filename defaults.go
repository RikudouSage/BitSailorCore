package bitwarden

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (receiver *client) provideDefaultsAndValidate() error {
	const defaultURL = "https://bitwarden.com"
	var normalized *urlConfig
	if receiver.identityURL == nil || receiver.apiURL == nil || receiver.notificationsURL == nil {
		normalized = normalizeBaseURL(lo.Must(url.Parse(defaultURL)))
	}

	if receiver.httpClient == nil {
		receiver.httpClient = http.DefaultClient
	}
	if receiver.identityURL == nil {
		receiver.identityURL = normalized.identityURL
	}
	if receiver.apiURL == nil {
		receiver.apiURL = normalized.apiURL
	}
	if receiver.notificationsURL == nil {
		receiver.notificationsURL = normalized.notificationsURL
	}

	if receiver.deviceID == uuid.Nil {
		return errors.New("device ID is required")
	}

	return nil
}
