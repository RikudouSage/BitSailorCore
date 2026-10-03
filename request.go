package bitwarden

import (
	"context"
	"net/url"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *auth) request[TResponse any](
	ctx context.Context,
	method string,
	url *url.URL,
	body any,
	session *result.Session,
) (TResponse, error) {
	return http.DoRequest[TResponse](ctx, receiver.httpClient, method, url, body, session, receiver.debugLogs)
}

func (receiver *vault) request[TResponse any](
	ctx context.Context,
	method string,
	url *url.URL,
	body any,
	session *result.Session,
) (TResponse, error) {
	return http.DoRequest[TResponse](ctx, receiver.httpClient, method, url, body, session, receiver.debugLogs)
}
