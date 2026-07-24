package bitwarden

import (
	"context"
	"fmt"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func (receiver *auth) UnlockSession(_ context.Context, session *result.Session, email, password string) error {
	if session.Encryption.UserKey != nil {
		return nil
	}

	if err := session.ValidateForUnlock(); err != nil {
		return fmt.Errorf("failed unlocking session: %w", err)
	}

	masterKey, err := crypto.DeriveMasterKey(email, password, session.Encryption.KDFType, &crypto.KDFConfig{
		Iterations:  session.Encryption.KDFIterations,
		Memory:      session.Encryption.KDFMemory,
		Parallelism: session.Encryption.KDFParallelism,
	})
	if err != nil {
		return fmt.Errorf("failed deriving master key: %w", err)
	}
	userKey, err := crypto.DecryptUserKey(*session.Encryption.EncryptedUserKey, masterKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt user key: %w", err)
	}

	session.Encryption.UserKey = userKey

	return nil
}
