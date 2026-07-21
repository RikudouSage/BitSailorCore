package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) Sync(ctx context.Context, session *result.Session) (*result.Sync, error) {
	if session == nil || session.Auth == nil {
		return nil, errors.New("session auth data is nil")
	}

	uri := new(*receiver.baseURL)
	uri.Path = "/sync"

	syncData, err := request[*result.Sync](
		ctx,
		receiver.httpClient,
		http.MethodGet,
		uri,
		nil,
		session,
	)
	if err != nil {
		return nil, fmt.Errorf("failed syncing vault: %w", err)
	}

	log.Println(syncData)

	return nil, nil
}
