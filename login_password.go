package bitwarden

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

var ErrTwoFactorRequired = fmt.Errorf("two factor authentication required")

func (receiver *auth) preLogin(ctx context.Context, email string) (*preLoginResponse, error) {
	uri := new(*receiver.identityURL)
	uri.Path = fmt.Sprintf("/identity/accounts/prelogin")

	resp, err := request[*preLoginResponse](ctx, receiver.httpClient, http.MethodPost, uri, &preLoginRequest{Email: email}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prelogin: %w", err)
	}

	return resp, nil
}

func (receiver *auth) LoginPassword(ctx context.Context, email, password string) (*result.Session, error) {
	preLogin, err := receiver.preLogin(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to prelogin: %w", err)
	}

	masterKey, err := crypto.DeriveMasterKey(email, password, preLogin.KDFType, &crypto.KDFConfig{
		Iterations:  preLogin.KDFIterations,
		Memory:      preLogin.KDFMemory,
		Parallelism: preLogin.KDFParallelism,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to derive master key: %w", err)
	}
	hash := crypto.DeriveMasterKeyHash(masterKey, password)

	requestData := &passwordLoginRequest{
		GrantType:        "password",
		Username:         email,
		Password:         hash,
		Scope:            "api offline_access",
		ClientID:         "cli",
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
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Bitwarden-Client-Version", internal.BitwardenVersion)

	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer internalHttp.DrainResponse(resp)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var twoFaResp twoFactorErrorResponse
	_ = json.Unmarshal(body, &twoFaResp)
	if len(twoFaResp.TwoFactorProviders) > 0 {
		return nil, ErrTwoFactorRequired
	}

	var token tokenResponse
	if err = json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	encryptedUserKey := token.GetUserKey()
	if encryptedUserKey == nil {
		return nil, fmt.Errorf("encrypted user key is nil")
	}

	userKey, err := crypto.DecryptUserKey(*encryptedUserKey, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt user key: %w", err)
	}

	return &result.Session{
		Auth: &result.AuthData{
			AccessToken:  token.AccessToken,
			ExpiresAt:    receiver.now().Add(time.Duration(token.ExpiresIn) * time.Second),
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		},
		Encryption: &result.EncryptionData{
			UserKey:             userKey,
			EncryptedPrivateKey: token.GetPrivateKey(),
			EncryptedUserKey:    token.GetUserKey(),

			KDFType:        token.KDFType,
			KDFIterations:  token.KDFIterations,
			KDFParallelism: token.KDFParallelism,
			KDFMemory:      token.KDFMemory,
		},
	}, nil
}
