package bitwarden

import (
	"fmt"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/dto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) getOrganizationKeys(session *result.Session) (map[uuid.UUID]dto.Key, error) {
	if session.Encryption.OrganizationKeys != nil {
		return session.Encryption.OrganizationKeys, nil
	}

	privateKeyBytes, err := crypto.DecryptBytes(*session.Encryption.EncryptedPrivateKey, session.Encryption.UserKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	privateKey, err := crypto.ParseRSAPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	decryptedOrgKeys := make(map[uuid.UUID]dto.Key)
	for orgId, orgKey := range receiver.vaultData.GetOrganizationKeys() {
		decryptedOrgKeys[orgId], err = crypto.DecryptRSAEncString(orgKey, privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt organization key: %w", err)
		}
	}

	return decryptedOrgKeys, nil
}
