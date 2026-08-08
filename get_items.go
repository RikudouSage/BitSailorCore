package bitwarden

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/types"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
	"golang.org/x/sync/errgroup"
)

func (receiver *vault) GetItems(ctx context.Context, session *result.Session) ([]*result.Item, error) {
	if receiver.vaultData == nil {
		return nil, ErrMissingVault
	}

	items := lo.Filter(receiver.vaultData.Items, func(item *result.Item, _ int) bool {
		return item.DeletedDate == nil
	})

	resultSlice := types.NewSyncSlice[*result.Item](len(items), len(items))

	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(20)

	for index, item := range items {
		wg.Go(func() error {
			newItem, err := receiver.DecryptItem(ctx, session, item)
			if err != nil {
				if !errors.Is(err, crypto.ErrInvalidEncryptedString) {
					return err
				}
				newItem = item.AsInvalidItem(err)
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
