package bitwarden

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	clone "github.com/huandu/go-clone/generic"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
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
	err := receiver.decryptSend(ctx, session, newSend)
	if err != nil {
		return nil, fmt.Errorf("failed decrypting the send item: %w", err)
	}

	return newSend, nil
}

func (receiver *vault) decryptSend(ctx context.Context, session *result.Session, send *result.Send) error {
	seed, err := crypto.DecryptBytes(send.Key, session.Encryption.UserKey)
	if err != nil {
		return fmt.Errorf("failed decrypting seed: %w", err)
	}

	sendKey, err := crypto.DeriveSendKey(seed)
	if err != nil {
		return fmt.Errorf("failed deriving send key: %w", err)
	}

	if err = receiver.decryptStruct(ctx, send, sendKey, []string{"Key"}); err != nil {
		return fmt.Errorf("failed decrypting send item: %w", err)
	}

	accessUri := fmt.Sprintf(
		"%s/#/send/%s/%s",
		receiver.sendURL,
		*send.AccessID,
		base64.RawURLEncoding.EncodeToString(seed),
	)

	send.AccessURL, err = url.Parse(accessUri)
	if err != nil {
		return fmt.Errorf("failed parsing access uri (%s): %w", accessUri, err)
	}

	return nil
}
