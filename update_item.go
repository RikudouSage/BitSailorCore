package bitwarden

import (
	"context"
	"fmt"
	"net/http"

	clone "github.com/huandu/go-clone/generic"
	"github.com/samber/lo"
	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *vault) UpdateItem(ctx context.Context, session *result.Session, item *result.Item) error {
	if receiver.vaultData == nil {
		return ErrMissingVault
	}

	resultItem := clone.Clone(item)
	key, err := receiver.getItemDecryptionKey(session, resultItem)
	if err != nil {
		return fmt.Errorf("failed fetching encryption key: %w", err)
	}
	err = receiver.encryptStruct(ctx, resultItem, key, []string{"Key"})
	if err != nil {
		return fmt.Errorf("failed encrypting struct: %w", err)
	}

	targetUri := internalHttp.UrlWithPath(receiver.apiURL, fmt.Sprintf("/ciphers/%s", item.ID))
	if err = receiver.auth.refreshIfNeeded(ctx, session); err != nil {
		return err
	}
	updatedItemEnc, err := receiver.request[*result.Item](ctx, http.MethodPut, targetUri, resultItem, session)
	if err != nil {
		return fmt.Errorf("failed updating the item: %w", err)
	}
	updatedItemDec, err := receiver.DecryptItem(ctx, session, updatedItemEnc)
	if err != nil {
		return fmt.Errorf("failed decrypting the item: %w", err)
	}
	*item = *updatedItemDec

	_, index, found := lo.FindIndexOf(receiver.vaultData.Items, func(item *result.Item) bool {
		return item.ID == updatedItemEnc.ID
	})
	if !found {
		return fmt.Errorf("error updating item with ID %s: %w", updatedItemEnc.ID, ErrItemNotFound)
	}
	receiver.vaultData.Items[index] = updatedItemEnc

	return nil
}
