package crypto

import (
	"testing"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
)

func TestDeriveSendPassword(t *testing.T) {
	t.Parallel()

	seed := dto.Key("0123456789abcdef0123456789abcdef")
	got := DeriveSendPassword("send-password", seed)
	const want = "Ujfz7e46sTPU6q0c2GB3AsghO0ozPEAdRcxzh5FNktM="
	if got != want {
		t.Fatalf("DeriveSendPassword() = %q, want %q", got, want)
	}
}
