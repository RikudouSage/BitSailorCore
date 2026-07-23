package bitwarden

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func TestDeleteSendDeletesRemoteSendAndRemovesItFromVaultData(t *testing.T) {
	t.Parallel()

	deleteID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	keepID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodDelete)
		}
		if r.URL.Path != "/sends/"+deleteID.String() {
			t.Fatalf("path = %s, want /sends/%s", r.URL.Path, deleteID)
		}
		if r.Header.Get("Authorization") != "Bearer access-token" {
			t.Fatalf("authorization header = %q, want Bearer access-token", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	vault := &vault{
		baseURL:    baseURL,
		httpClient: server.Client(),
		vaultData: &result.VaultData{
			Sends: []*result.Send{
				{ID: keepID},
				{ID: deleteID},
			},
		},
	}
	session := &result.Session{
		Auth: &result.AuthData{
			AccessToken: "access-token",
			TokenType:   "Bearer",
		},
	}

	if err = vault.DeleteSend(context.Background(), session, deleteID); err != nil {
		t.Fatalf("DeleteSend() returned error: %v", err)
	}

	if len(vault.vaultData.Sends) != 1 {
		t.Fatalf("len(vaultData.Sends) = %d, want 1", len(vault.vaultData.Sends))
	}
	if vault.vaultData.Sends[0].ID != keepID {
		t.Fatalf("remaining send ID = %s, want %s", vault.vaultData.Sends[0].ID, keepID)
	}
}
