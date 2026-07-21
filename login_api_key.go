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

func (receiver *auth) LoginApiKey(ctx context.Context, clientID, clientSecret string) (*result.Session, error) {
	requestData := &apiKeyLoginRequest{
		GrantType:        "client_credentials",
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		Scope:            "api",
		DeviceType:       deviceTypeLinuxCLI,
		DeviceIdentifier: receiver.deviceID.String(),
		DeviceName:       internal.DeviceName,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		receiver.getTokenURL().String(),
		strings.NewReader(lo.Must(query.Values(requestData)).Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Bitwarden-Client-Version", internal.BitwardenVersion)

	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed sending request: %w", err)
	}
	defer internalHttp.DrainResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var token tokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &result.Session{
		Auth: &result.AuthData{
			AccessToken:  token.AccessToken,
			ExpiresAt:    receiver.now().Add(time.Duration(token.ExpiresIn) * time.Second),
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		},
		Encryption: &result.EncryptionData{
			EncryptedPrivateKey: token.GetPrivateKey(),
			EncryptedUserKey:    token.GetUserKey(),

			KDFType:        token.KDFType,
			KDFIterations:  token.KDFIterations,
			KDFParallelism: token.KDFParallelism,
			KDFMemory:      token.KDFMemory,
		},
	}, nil
}
