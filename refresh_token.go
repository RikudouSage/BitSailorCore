package bitwarden

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal"
	internalHttp "go.chrastecky.dev/bitwarden-client/bitwarden/internal/http"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *auth) RefreshToken(ctx context.Context, session *result.Session) error {
	requestData := &refreshLoginRequest{
		GrantType:    "refresh_token",
		ClientID:     "cli",
		ClientSecret: session.Auth.RefreshToken,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		receiver.getTokenURL().String(),
		strings.NewReader(lo.Must(query.Values(requestData)).Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Bitwarden-Client-Version", internal.BitwardenVersion)

	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed sending request: %w", err)
	}
	defer internalHttp.DrainResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var token tokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	session.Auth.AccessToken = token.AccessToken
	session.Auth.ExpiresAt = receiver.now().Add(time.Duration(token.ExpiresIn) * time.Second)
	session.Auth.RefreshToken = token.RefreshToken
	session.Auth.TokenType = token.TokenType

	return nil
}
