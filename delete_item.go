package bitwarden

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/samber/lo"
	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *vault) DeleteItem(ctx context.Context, session *result.Session, itemID uuid.UUID) error {
	if receiver.vaultData == nil {
		return ErrMissingVault
	}

	targetUri := internalHttp.UrlWithPath(receiver.apiURL, fmt.Sprintf("/ciphers/%s/delete", itemID))

	if err := receiver.auth.refreshIfNeeded(ctx, session); err != nil {
		return err
	}
	_, err := receiver.request[any](ctx, http.MethodPut, targetUri, nil, session)
	if err != nil {
		return fmt.Errorf("failed deleting item: %w", err)
	}

	receiver.vaultData.Items = lo.Filter(receiver.vaultData.Items, func(item *result.Item, _ int) bool {
		return item.ID != itemID
	})
	return nil
}
