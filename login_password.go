package bitwarden

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitsailor-core/bitwarden/dto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	internalHttp "go.chrastecky.dev/bitsailor-core/bitwarden/internal/http"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

var ErrTwoFactorRequired = fmt.Errorf("two factor authentication required")
var ErrUnsupportedTwoFactorRequired = fmt.Errorf("unsupported two factor authentication required")

func (receiver *auth) preLogin(ctx context.Context, email string) (*preLoginResponse, error) {
	uri := internalHttp.UrlWithPath(receiver.identityURL, "/identity/accounts/prelogin")

	resp, err := receiver.request[*preLoginResponse](ctx, http.MethodPost, uri, &preLoginRequest{Email: email}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prelogin: %w", err)
	}

	return resp, nil
}

func (receiver *auth) LoginPassword(ctx context.Context, email, password string, tfa *dto.TFAConfig) (*result.Session, error) {
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

	if tfa != nil {
		provider, err := receiver.providers.FindByKind(tfa.Kind)
		if err != nil {
			return nil, fmt.Errorf("failed to find provider: %w", err)
		}

		if err = provider.ModifyRequest(requestData); err != nil {
			return nil, fmt.Errorf("failed to modify request using totp provider: %w", err)
		}

		requestData.TwoFactorToken = new(tfa.Code)
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
	if len(twoFaResp.TwoFactorProviders2) > 0 {
		appSupportedKinds := receiver.providers.GetSupportedKinds()
		accountSupportedKinds, err := lo.MapErr(lo.Keys(twoFaResp.TwoFactorProviders2), func(item string, _ int) (dto.TFAKind, error) {
			parsed, err := strconv.Atoi(item)
			if err != nil {
				return 0, fmt.Errorf("failed parsing '%s' as a provider type: %w", item, err)
			}

			return dto.TFAKind(parsed), nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to parse two-factor authentication provider kinds: %w", err)
		}

		intersection := lo.Intersect(appSupportedKinds, accountSupportedKinds)
		if len(intersection) == 0 {
			return nil, fmt.Errorf(
				"%w: supported methods: %s",
				ErrUnsupportedTwoFactorRequired,
				strings.Join(receiver.providers.GetSupportedNames(), ", "),
			)
		}

		supportedKindsStrs := lo.Map(intersection, func(item dto.TFAKind, _ int) string {
			return strconv.Itoa(int(item))
		})
		return nil, fmt.Errorf(
			"%w: supported kinds: %s|",
			ErrTwoFactorRequired,
			strings.Join(supportedKindsStrs, ", "),
		)
	}

	if resp.StatusCode != http.StatusOK {
		if receiver.debugLogs {
			return nil, fmt.Errorf("unexpected status code: %d (%s), body: %s", resp.StatusCode, resp.Status, body)
		}
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
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
