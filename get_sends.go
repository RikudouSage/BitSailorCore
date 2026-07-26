package bitwarden

import (
	"context"

	clone "github.com/huandu/go-clone/generic"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/types"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
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
			err := receiver.decryptSend(ctx, session, newSend)
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
