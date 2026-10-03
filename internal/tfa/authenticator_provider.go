package tfa

import (
	"context"

	"go.chrastecky.dev/bitsailor-core/bitwarden/dto"
)

type AuthenticatorProvider interface {
	Provider
}

type authenticatorProvider struct {
}

func NewAuthenticatorProvider() AuthenticatorProvider {
	return &authenticatorProvider{}
}

func (receiver *authenticatorProvider) Initialize(_ context.Context, _, _ string) error {
	return nil
}

func (receiver *authenticatorProvider) Kind() dto.TFAKind {
	return dto.TFAKindAuthenticator
}

func (receiver *authenticatorProvider) ModifyRequest(request Request) error {
	request.SetProvider(new(receiver.Kind()))
	request.SetRemember(new(false))

	return nil
}

func (receiver *authenticatorProvider) Name() string {
	return "authenticator"
}
