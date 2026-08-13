package result

import (
	"testing"

	"github.com/google/uuid"
)

func TestVaultDataGetOrganizationKeys(t *testing.T) {
	t.Parallel()

	firstID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	secondID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	vaultData := VaultData{
		Profile: &Profile{
			Organizations: []*ProfileOrganization{
				{ID: firstID, Key: "first-key"},
				{ID: secondID, Key: "second-key"},
			},
		},
	}

	got := vaultData.GetOrganizationKeys()
	if len(got) != 2 {
		t.Fatalf("len(keys) = %d, want 2", len(got))
	}
	if got[firstID] != "first-key" {
		t.Fatalf("keys[firstID] = %q, want first-key", got[firstID])
	}
	if got[secondID] != "second-key" {
		t.Fatalf("keys[secondID] = %q, want second-key", got[secondID])
	}
}
