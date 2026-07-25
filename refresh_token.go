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
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal"
	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *auth) RefreshToken(ctx context.Context, session *result.Session) error {
	if session.Auth.RefreshToken == "" && session.Auth.ClientID != nil && session.Auth.ClientSecret != nil {
		newSession, err := receiver.LoginApiKey(ctx, *session.Auth.ClientID, *session.Auth.ClientSecret)
		if err != nil {
			return fmt.Errorf("failed refreshing through new login via client id / client secret: %w", err)
		}

		session.Auth.TokenType = newSession.Auth.TokenType
		session.Auth.AccessToken = newSession.Auth.AccessToken
		session.Auth.RefreshToken = newSession.Auth.RefreshToken
		session.Auth.ExpiresAt = newSession.Auth.ExpiresAt
		return nil
	}

	requestData := &refreshLoginRequest{
		GrantType:    "refresh_token",
		ClientID:     "cli",
		RefreshToken: session.Auth.RefreshToken,
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

func (receiver *auth) refreshIfNeeded(ctx context.Context, session *result.Session) error {
	if session == nil || session.Auth == nil {
		return nil
	}
	if session.Auth.ExpiresAt.After(receiver.now()) {
		return nil
	}

	err := receiver.RefreshToken(ctx, session)
	if err != nil {
		return fmt.Errorf("failed refreshing token: %w", err)
	}

	return nil
}
