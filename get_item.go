package bitwarden

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) GetItem(ctx context.Context, session *result.Session, itemID uuid.UUID) (*result.Item, error) {
	if receiver.vaultData == nil {
		return nil, ErrMissingVault
	}

	originalItem, found := lo.Find(receiver.vaultData.Items, func(item *result.Item) bool {
		return item.ID == itemID
	})
	if !found {
		return nil, fmt.Errorf("error getting item with ID %s: %w", itemID, ErrItemNotFound)
	}

	return receiver.DecryptItem(ctx, session, originalItem)
}
