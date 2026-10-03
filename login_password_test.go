package bitwarden

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/dto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/tfa"
)

func TestLoginPasswordUsesSelectedTwoFactorProvider(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/accounts/prelogin":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"kdf":0,"kdfIterations":1}`))
		case "/identity/connect/token":
			if err := r.ParseForm(); err != nil {
				t.Fatalf("ParseForm() returned error: %v", err)
			}
			if got := r.Form.Get("twoFactorProvider"); got != "0" {
				t.Fatalf("twoFactorProvider = %q, want 0", got)
			}
			if got := r.Form.Get("twoFactorToken"); got != "123456" {
				t.Fatalf("twoFactorToken = %q, want 123456", got)
			}
			if got := r.Form.Get("twoFactorRemember"); got != "0" {
				t.Fatalf("twoFactorRemember = %q, want 0", got)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"TwoFactorProviders2":{"0":{}}}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	auth := testAuth(t, server)
	_, err := auth.LoginPassword(context.Background(), "person@example.test", "password", &dto.TFAConfig{
		Kind: dto.TFAKindAuthenticator,
		Code: "123456",
	})
	if !errors.Is(err, ErrTwoFactorRequired) {
		t.Fatalf("LoginPassword() error = %v, want ErrTwoFactorRequired", err)
	}
}

func TestLoginPasswordRejectsAccountWithoutSupportedProvider(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/accounts/prelogin":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"kdf":0,"kdfIterations":1}`))
		case "/identity/connect/token":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"TwoFactorProviders2":{"1":{}}}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	auth := testAuth(t, server)
	_, err := auth.LoginPassword(context.Background(), "person@example.test", "password", nil)
	if !errors.Is(err, ErrUnsupportedTwoFactorRequired) {
		t.Fatalf("LoginPassword() error = %v, want ErrUnsupportedTwoFactorRequired", err)
	}
	if !strings.Contains(err.Error(), "authenticator") {
		t.Fatalf("LoginPassword() error = %q, want supported provider name", err)
	}
}

func TestInitializeTOTPProviderSendsEmailCode(t *testing.T) {
	t.Parallel()

	deviceID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/accounts/prelogin":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"kdf":0,"kdfIterations":1}`))
		case "/two-factor/send-email-login":
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
			}
			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}

			var body struct {
				Email              string `json:"email"`
				MasterPasswordHash string `json:"masterPasswordHash"`
				DeviceIdentifier   string `json:"deviceIdentifier"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("Decode() returned error: %v", err)
			}
			if body.Email != "person@example.test" {
				t.Fatalf("email = %q, want person@example.test", body.Email)
			}
			if body.MasterPasswordHash == "" {
				t.Fatal("masterPasswordHash is empty")
			}
			if body.DeviceIdentifier != deviceID.String() {
				t.Fatalf("deviceIdentifier = %q, want %s", body.DeviceIdentifier, deviceID)
			}

			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	auth := testAuthWithProviders(t, server, deviceID, tfa.Providers{
		tfa.NewAuthenticatorProvider(),
		tfa.NewEmailProvider(server.Client(), mustParseURL(t, server.URL), deviceID, false),
	})
	if err := auth.InitializeTOTPProvider(context.Background(), "person@example.test", "password", dto.TFAKindEmail); err != nil {
		t.Fatalf("InitializeTOTPProvider() returned error: %v", err)
	}
}

func TestLoginPasswordUsesEmailTwoFactorProvider(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/accounts/prelogin":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"kdf":0,"kdfIterations":1}`))
		case "/identity/connect/token":
			if err := r.ParseForm(); err != nil {
				t.Fatalf("ParseForm() returned error: %v", err)
			}
			if got := r.Form.Get("twoFactorProvider"); got != "1" {
				t.Fatalf("twoFactorProvider = %q, want 1", got)
			}
			if got := r.Form.Get("twoFactorToken"); got != "654321" {
				t.Fatalf("twoFactorToken = %q, want 654321", got)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"TwoFactorProviders2":{"1":{}}}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	auth := testAuthWithProviders(t, server, uuid.MustParse("11111111-1111-1111-1111-111111111111"), tfa.Providers{
		tfa.NewAuthenticatorProvider(),
		tfa.NewEmailProvider(server.Client(), mustParseURL(t, server.URL), uuid.MustParse("11111111-1111-1111-1111-111111111111"), false),
	})
	_, err := auth.LoginPassword(context.Background(), "person@example.test", "password", &dto.TFAConfig{
		Kind: dto.TFAKindEmail,
		Code: "654321",
	})
	if !errors.Is(err, ErrTwoFactorRequired) {
		t.Fatalf("LoginPassword() error = %v, want ErrTwoFactorRequired", err)
	}
}

func testAuth(t *testing.T, server *httptest.Server) *auth {
	t.Helper()
	return testAuthWithProviders(t, server, uuid.MustParse("11111111-1111-1111-1111-111111111111"), tfa.Providers{tfa.NewAuthenticatorProvider()})
}

func testAuthWithProviders(t *testing.T, server *httptest.Server, deviceID uuid.UUID, providers tfa.Providers) *auth {
	t.Helper()

	baseURL := mustParseURL(t, server.URL)

	return newAuth(
		baseURL,
		baseURL,
		server.Client(),
		deviceID,
		false,
		providers,
	)
}
