package bitwarden

import (
	"context"
	"fmt"

	clone "github.com/huandu/go-clone/generic"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/dto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/types"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
	"golang.org/x/sync/errgroup"
)

func (receiver *vault) GetSends(ctx context.Context, session *result.Session) ([]*result.Send, error) {
	if receiver.vaultData == nil {
		return nil, ErrMissingVault
	}

	resultSlice := types.NewSyncSlice[*result.Send](len(receiver.vaultData.Sends), len(receiver.vaultData.Sends))

	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(20)

	for index, send := range receiver.vaultData.Sends {
		wg.Go(func() error {
			newSend := clone.Clone(send)
			key, err := receiver.getSendDecryptionKey(newSend, session.Encryption.UserKey)
			if err != nil {
				return err
			}
			err = receiver.decryptStruct(ctx, newSend, key, []string{"Key"})
			if err != nil {
				return err
			}
			if err = resultSlice.Insert(index, newSend); err != nil {
				return err
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return resultSlice.ToSlice(), nil
}

func (*vault) getSendDecryptionKey(send *result.Send, userKey dto.Key) (dto.Key, error) {
	seed, err := crypto.DecryptBytes(send.Key, userKey)
	if err != nil {
		return nil, fmt.Errorf("failed decrypting seed: %w", err)
	}

	return crypto.DeriveSendKey(seed)
}
