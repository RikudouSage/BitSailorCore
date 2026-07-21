package bitwarden

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (receiver *client) provideDefaultsAndValidate() error {
	if receiver.httpClient == nil {
		receiver.httpClient = http.DefaultClient
	}
	if receiver.baseURL == nil {
		receiver.baseURL = lo.Must(url.Parse("https://vault.bitwarden.com"))
	}

	if receiver.deviceID == uuid.Nil {
		return errors.New("device ID is required")
	}

	return nil
}
