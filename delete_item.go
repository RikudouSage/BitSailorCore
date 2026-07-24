package bitwarden

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *vault) DeleteItem(ctx context.Context, session *result.Session, itemID uuid.UUID) error {
	if receiver.vaultData == nil {
		return ErrMissingVault
	}

	targetUri := new(*receiver.baseURL)
	targetUri.Path = fmt.Sprintf("/ciphers/%s/delete", itemID)

	_, err := request[any](ctx, receiver.httpClient, http.MethodPut, targetUri, nil, session)
	if err != nil {
		return fmt.Errorf("failed deleting item: %w", err)
	}

	receiver.vaultData.Items = lo.Filter(receiver.vaultData.Items, func(item *result.Item, _ int) bool {
		return item.ID != itemID
	})
	return nil
}
