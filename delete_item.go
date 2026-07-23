package bitwarden

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) DeleteItem(ctx context.Context, session *result.Session, itemID uuid.UUID) error {
	targetUri := new(*receiver.baseURL)
	targetUri.Path = fmt.Sprintf("/ciphers/%s/delete", itemID)

	_, err := request[any](ctx, receiver.httpClient, http.MethodPut, targetUri, nil, session)
	if err != nil {
		return fmt.Errorf("failed deleting item: %w", err)
	}

	return nil
}
