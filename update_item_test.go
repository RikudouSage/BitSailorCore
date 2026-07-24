package bitwarden

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func TestUpdateItemUpdatesRemoteItemAndReplacesVaultData(t *testing.T) {
	t.Parallel()

	itemID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	keepID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	userKey := dto.Key([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	originalEncryptedName, err := crypto.EncryptString("Original item", userKey)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}
	updatedEncryptedName, err := crypto.EncryptString("Updated item", userKey)
	if err != nil {
		t.Fatalf("EncryptString() returned error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/ciphers/"+itemID.String() {
			t.Fatalf("path = %s, want /ciphers/%s", r.URL.Path, itemID)
		}
		if r.Header.Get("Authorization") != "Bearer access-token" {
			t.Fatalf("authorization header = %q, want Bearer access-token", r.Header.Get("Authorization"))
		}

		var requestItem result.Item
		if err := json.NewDecoder(r.Body).Decode(&requestItem); err != nil {
			t.Fatalf("Decode() returned error: %v", err)
		}
		if requestItem.ID != itemID {
			t.Fatalf("request item ID = %s, want %s", requestItem.ID, itemID)
		}
		if requestItem.Name == "Updated item" {
			t.Fatal("request item name was not encrypted")
		}
		decryptedName, err := crypto.DecryptString(requestItem.Name, userKey)
		if err != nil {
			t.Fatalf("DecryptString() returned error: %v", err)
		}
		if decryptedName != "Updated item" {
			t.Fatalf("decrypted request item name = %q, want Updated item", decryptedName)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&result.Item{
			ID:         itemID,
			Type:       result.ItemTypeSecureNote,
			Name:       updatedEncryptedName,
			SecureNote: &result.ItemSecureNote{},
		})
		if err != nil {
			t.Fatalf("Encode() returned error: %v", err)
		}
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
			Items: []*result.Item{
				{ID: keepID},
				{ID: itemID, Type: result.ItemTypeSecureNote, Name: originalEncryptedName, SecureNote: &result.ItemSecureNote{}},
			},
		},
	}
	session := &result.Session{
		Auth: &result.AuthData{
			AccessToken: "access-token",
			TokenType:   "Bearer",
		},
		Encryption: &result.EncryptionData{
			UserKey: userKey,
		},
	}
	item := &result.Item{
		ID:         itemID,
		Type:       result.ItemTypeSecureNote,
		Name:       "Updated item",
		SecureNote: &result.ItemSecureNote{},
	}

	if err = vault.UpdateItem(context.Background(), session, item); err != nil {
		t.Fatalf("UpdateItem() returned error: %v", err)
	}

	if item.Name != "Updated item" {
		t.Fatalf("item.Name = %q, want Updated item", item.Name)
	}
	if len(vault.vaultData.Items) != 2 {
		t.Fatalf("len(vaultData.Items) = %d, want 2", len(vault.vaultData.Items))
	}
	if vault.vaultData.Items[0].ID != keepID {
		t.Fatalf("first cached item ID = %s, want %s", vault.vaultData.Items[0].ID, keepID)
	}
	if vault.vaultData.Items[1].ID != itemID {
		t.Fatalf("updated cached item ID = %s, want %s", vault.vaultData.Items[1].ID, itemID)
	}
	if vault.vaultData.Items[1].Name == "Updated item" {
		t.Fatal("cached item name was not encrypted")
	}
	decryptedCachedName, err := crypto.DecryptString(vault.vaultData.Items[1].Name, userKey)
	if err != nil {
		t.Fatalf("DecryptString() returned error: %v", err)
	}
	if decryptedCachedName != "Updated item" {
		t.Fatalf("decrypted cached item name = %q, want Updated item", decryptedCachedName)
	}
}
