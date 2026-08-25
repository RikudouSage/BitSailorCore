package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/types"
	"golang.org/x/sync/errgroup"
)

type Icons interface {
	GetIcon(ctx context.Context, hostname string) ([]byte, error)
	GetIcons(ctx context.Context, hostnames []string) (map[string][]byte, error)
}

type icons struct {
	httpClient *http.Client
	debugLogs  bool
	uri        *url.URL
}

func newIcons(
	httpClient *http.Client,
	debugLogs bool,
	uri *url.URL,
) *icons {
	return &icons{
		httpClient: httpClient,
		debugLogs:  debugLogs,
		uri:        uri,
	}
}

func (receiver *icons) request(ctx context.Context, uri *url.URL) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Accept", "image/png")

	resp, err := receiver.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed sending icon request: %w", err)
	}
	defer internalHttp.DrainResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed sending icon request: unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	return body, nil
}

func (receiver *icons) GetIcon(ctx context.Context, hostname string) ([]byte, error) {
	return receiver.request(ctx, receiver.url(hostname))
}

func (receiver *icons) url(hostname string) *url.URL {
	uri := new(*receiver.uri)
	basePath := strings.TrimRight(uri.Path, "/")
	escapedBasePath := strings.TrimRight(uri.EscapedPath(), "/")

	uri.Path = fmt.Sprintf("%s/%s/icon.png", basePath, hostname)
	uri.RawPath = fmt.Sprintf("%s/%s/icon.png", escapedBasePath, url.PathEscape(hostname))

	return uri
}

func (receiver *icons) GetIcons(ctx context.Context, hostnames []string) (map[string][]byte, error) {
	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(5)

	results := types.NewSyncMap[string, []byte](len(hostnames))
	errs := types.NewSyncSlice[error](0, 0)

	for _, hostname := range hostnames {
		wg.Go(func() error {
			icon, err := receiver.GetIcon(ctx, hostname)
			if err != nil {
				errs.Append(fmt.Errorf("failed getting icon for '%s': %w", hostname, err))
				return nil
			}

			results.Insert(hostname, icon)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return results.ToMap(), errors.Join(errs.ToSlice()...)
}
