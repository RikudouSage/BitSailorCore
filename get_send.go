package bitwarden

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	clone "github.com/huandu/go-clone/generic"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) GetSend(ctx context.Context, session *result.Session, itemID uuid.UUID) (*result.Send, error) {
	if receiver.vaultData == nil {
		return nil, ErrMissingVault
	}

	originalItem, found := lo.Find(receiver.vaultData.Sends, func(item *result.Send) bool {
		return item.ID == itemID
	})
	if !found {
		return nil, fmt.Errorf("error getting send with ID %s: %w", itemID, ErrItemNotFound)
	}

	newSend := clone.Clone(originalItem)
	key, err := receiver.getSendDecryptionKey(newSend, session.Encryption.UserKey)
	if err != nil {
		return nil, fmt.Errorf("failed getting decryption key: %w", err)
	}
	err = receiver.decryptStruct(ctx, newSend, key, []string{"Key"})
	if err != nil {
		return nil, fmt.Errorf("failed decrypting the send item: %w", err)
	}

	return newSend, nil
}
