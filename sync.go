package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *vault) Sync(ctx context.Context, session *result.Session) (Vault, error) {
	if session == nil || session.Auth == nil {
		return nil, errors.New("session auth data is nil")
	}

	if err := receiver.auth.refreshIfNeeded(ctx, session); err != nil {
		return nil, err
	}

	uri := internalHttp.UrlWithPath(receiver.apiURL, "/sync")

	vaultData, err := receiver.request[*result.VaultData](
		ctx,
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
