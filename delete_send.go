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

	targetUri := new(*receiver.baseURL)
	targetUri.Path = fmt.Sprintf("/sends/%s", sendID)

	_, err := request[any](ctx, receiver.httpClient, http.MethodDelete, targetUri, nil, session)
	if err != nil {
		return fmt.Errorf("failed deleting send: %w", err)
	}

	receiver.vaultData.Sends = lo.Filter(receiver.vaultData.Sends, func(send *result.Send, _ int) bool {
		return send.ID != sendID
	})
	return nil
}
