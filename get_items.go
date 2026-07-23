package bitwarden

import (
	"context"

	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/types"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
	"golang.org/x/sync/errgroup"
)

func (receiver *vault) GetItems(ctx context.Context, session *result.Session) ([]*result.Item, error) {
	if receiver.vaultData == nil {
		return nil, ErrMissingVault
	}

	resultSlice := types.NewSyncSlice[*result.Item](len(receiver.vaultData.Items), len(receiver.vaultData.Items))

	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(20)

	for index, item := range receiver.vaultData.Items {
		wg.Go(func() error {
			newItem, err := receiver.DecryptItem(ctx, session, item)
			if err != nil {
				return err
			}
			if err = resultSlice.Insert(index, newItem); err != nil {
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
