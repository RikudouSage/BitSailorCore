package bitwarden

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *auth) FetchAuthRequest(ctx context.Context, session *result.Session, id uuid.UUID) (*result.AuthRequest, error) {
	req, err := request[*result.AuthRequest](
		ctx,
		receiver.httpClient,
		http.MethodGet,
		urlWithPath(receiver.apiURL, fmt.Sprintf("/auth-requests/%s", id)),
		nil,
		session,
		receiver.debugLogs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed getting auth request: %w", err)
	}

	return req, nil
}

func (receiver *auth) RespondToAuthRequest(ctx context.Context, session *result.Session, authRequest *result.AuthRequest, approved bool) error {
	if session.Encryption.UserKey == nil {
		return ErrLockedSession
	}

	update := &authRequestUpdateRequest{
		DeviceIdentifier: receiver.deviceID,
		RequestApproved:  approved,
	}
	if approved {
		key, err := crypto.EncryptRSAEncBytes(session.Encryption.UserKey, authRequest.PublicKey)
		if err != nil {
			return fmt.Errorf("failed encrypting user key with request's public key: %w", err)
		}

		update.Key = &key
	}

	resp, err := request[*result.AuthRequest](
		ctx,
		receiver.httpClient,
		http.MethodPut,
		urlWithPath(receiver.apiURL, fmt.Sprintf("/auth-requests/%s", authRequest.ID)),
		update,
		session,
		receiver.debugLogs,
	)
	if err != nil {
		return fmt.Errorf("failed responding to auth request: %w", err)
	}

	if resp.RequestApproved == nil || *resp.RequestApproved != approved {
		return fmt.Errorf("auth request response did not confirm requested approval state")
	}

	return nil
}
