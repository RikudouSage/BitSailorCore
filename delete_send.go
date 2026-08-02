package bitwarden

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *vault) DeleteSend(ctx context.Context, session *result.Session, sendID uuid.UUID) error {
	if receiver.vaultData == nil {
		return ErrMissingVault
	}

	targetUri := urlWithPath(receiver.apiURL, fmt.Sprintf("/sends/%s", sendID))

	if err := receiver.auth.refreshIfNeeded(ctx, session); err != nil {
		return err
	}
	_, err := request[any](ctx, receiver.httpClient, http.MethodDelete, targetUri, nil, session, false)
	if err != nil {
		return fmt.Errorf("failed deleting send: %w", err)
	}

	receiver.vaultData.Sends = lo.Filter(receiver.vaultData.Sends, func(send *result.Send, _ int) bool {
		return send.ID != sendID
	})
	return nil
}
