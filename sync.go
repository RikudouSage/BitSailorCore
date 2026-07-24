package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *vault) Sync(ctx context.Context, session *result.Session) (Vault, error) {
	if session == nil || session.Auth == nil {
		return nil, errors.New("session auth data is nil")
	}

	uri := new(*receiver.baseURL)
	uri.Path = "/sync"

	vaultData, err := request[*result.VaultData](
		ctx,
		receiver.httpClient,
		http.MethodGet,
		uri,
		nil,
		session,
	)
	if err != nil {
		return nil, fmt.Errorf("failed syncing vault: %w", err)
	}

	return receiver.WithVaultData(vaultData), nil
}
