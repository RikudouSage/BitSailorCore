package bitwarden

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func testAuth(t *testing.T, server *httptest.Server) *auth {
	t.Helper()

	baseURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	return newAuth(
		baseURL,
		baseURL,
		server.Client(),
		uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		false,
		tfa.Providers{tfa.NewAuthenticatorProvider()},
	)
}
