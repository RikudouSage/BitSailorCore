package bitwarden

import (
	"context"
	"fmt"

	"go.chrastecky.dev/bitsailor-core/bitwarden/dto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
)

func (receiver *auth) InitializeTOTPProvider(ctx context.Context, email, password string, kind dto.TFAKind) error {
	provider, err := receiver.providers.FindByKind(kind)
	if err != nil {
		return fmt.Errorf("failed finding totp provider: %w", err)
	}

	preLogin, err := receiver.preLogin(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to prelogin: %w", err)
	}

	masterKey, err := crypto.DeriveMasterKey(email, password, preLogin.KDFType, &crypto.KDFConfig{
		Iterations:  preLogin.KDFIterations,
		Memory:      preLogin.KDFMemory,
		Parallelism: preLogin.KDFParallelism,
	})
	if err != nil {
		return fmt.Errorf("failed to derive master key: %w", err)
	}
	hash := crypto.DeriveMasterKeyHash(masterKey, password)

	if err = provider.Initialize(ctx, email, hash); err != nil {
		return fmt.Errorf("failed to initialize totp provider: %w", err)
	}

	return nil
}
