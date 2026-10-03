package tfa

import (
	"context"
	"fmt"

	"go.chrastecky.dev/bitsailor-core/bitwarden/dto"
)

type Provider interface {
	Kind() dto.TFAKind
	Name() string
	Initialize(ctx context.Context, email, hash string) error
	ModifyRequest(Request) error
}

type Providers []Provider

func (receiver Providers) FindByKind(kind dto.TFAKind) (Provider, error) {
	for _, provider := range receiver {
		if provider.Kind() == kind {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("provider with kind %d not found", kind)
}

func (receiver Providers) GetSupportedKinds() []dto.TFAKind {
	result := make([]dto.TFAKind, len(receiver))

	for i, provider := range receiver {
		result[i] = provider.Kind()
	}

	return result
}

func (receiver Providers) GetSupportedNames() []string {
	result := make([]string, len(receiver))

	for i, provider := range receiver {
		result[i] = provider.Name()
	}

	return result
}
