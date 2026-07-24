package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	clone "github.com/huandu/go-clone/generic"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) UpdateItem(ctx context.Context, session *result.Session, item *result.Item) error {
	if receiver.vaultData == nil {
		return ErrMissingVault
	}

	if item.OrganizationID != uuid.Nil {
		return errors.New("updating items inside organizations is not supported yet")
	}

	resultItem := clone.Clone(item)
	err := receiver.encryptStruct(ctx, resultItem, session.Encryption.UserKey, nil)
	if err != nil {
		return fmt.Errorf("failed encrypting struct: %w", err)
	}

	targetUri := new(*receiver.baseURL)
	targetUri.Path = fmt.Sprintf("/ciphers/%s", item.ID)
	updatedItemEnc, err := request[*result.Item](ctx, receiver.httpClient, http.MethodPut, targetUri, resultItem, session)
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
