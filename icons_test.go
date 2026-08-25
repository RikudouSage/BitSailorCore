package bitwarden

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetIconRequestsPngForHostname(t *testing.T) {
	t.Parallel()

	const icon = "png-bytes"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/example.com/icon.png" {
			t.Fatalf("path = %s, want /example.com/icon.png", r.URL.Path)
		}
		if accept := r.Header.Get("Accept"); accept != "image/png" {
			t.Fatalf("Accept = %q, want image/png", accept)
		}

		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte(icon))
	}))
	t.Cleanup(server.Close)

	iconsService := newIcons(server.Client(), false, mustParseURL(t, server.URL))
	actual, err := iconsService.GetIcon(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("GetIcon() returned error: %v", err)
	}
	if string(actual) != icon {
		t.Fatalf("GetIcon() = %q, want %q", actual, icon)
	}
}

func TestGetIconEscapesHostnamePathSegment(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/example.com%2Fextra/icon.png" {
			t.Fatalf("escaped path = %s, want /example.com%%2Fextra/icon.png", r.URL.EscapedPath())
		}

		_, _ = w.Write([]byte("icon"))
	}))
	t.Cleanup(server.Close)

	iconsService := newIcons(server.Client(), false, mustParseURL(t, server.URL))
	if _, err := iconsService.GetIcon(context.Background(), "example.com/extra"); err != nil {
		t.Fatalf("GetIcon() returned error: %v", err)
	}
}

func TestGetIconReturnsErrorForUnexpectedStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	iconsService := newIcons(server.Client(), false, mustParseURL(t, server.URL))
	_, err := iconsService.GetIcon(context.Background(), "missing.example")
	if err == nil {
		t.Fatal("GetIcon() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "unexpected status code: 404") {
		t.Fatalf("GetIcon() error = %v, want status code error", err)
	}
}

func TestGetIconsReturnsPartialResultsAndJoinedErrors(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok.example/icon.png":
			_, _ = w.Write([]byte("ok-icon"))
		case "/missing.example/icon.png":
			http.Error(w, "not found", http.StatusNotFound)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	iconsService := newIcons(server.Client(), false, mustParseURL(t, server.URL))
	results, err := iconsService.GetIcons(context.Background(), []string{"ok.example", "missing.example"})
	if err == nil {
		t.Fatal("GetIcons() error = nil, want partial error")
	}
	if !strings.Contains(err.Error(), "failed getting icon for 'missing.example'") {
		t.Fatalf("GetIcons() error = %v, want hostname context", err)
	}
	if !strings.Contains(err.Error(), "unexpected status code: 404") {
		t.Fatalf("GetIcons() error = %v, want status code context", err)
	}

	if string(results["ok.example"]) != "ok-icon" {
		t.Fatalf("ok.example icon = %q, want ok-icon", results["ok.example"])
	}
	if _, ok := results["missing.example"]; ok {
		t.Fatal("missing.example was present in results, want omitted failed icon")
	}
}

func TestGetIconsReturnsNilErrorWhenAllIconsSucceed(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Path))
	}))
	t.Cleanup(server.Close)

	iconsService := newIcons(server.Client(), false, mustParseURL(t, server.URL))
	results, err := iconsService.GetIcons(context.Background(), []string{"one.example", "two.example"})
	if err != nil {
		t.Fatalf("GetIcons() returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
}

func TestGetIconsJoinedErrorMatchesUnderlyingStatusError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	iconsService := newIcons(server.Client(), false, mustParseURL(t, server.URL))
	_, err := iconsService.GetIcons(context.Background(), []string{"one.example"})
	if err == nil {
		t.Fatal("GetIcons() error = nil, want error")
	}

	expected := errors.New("unexpected status code: 503")
	if !strings.Contains(err.Error(), expected.Error()) {
		t.Fatalf("GetIcons() error = %v, want %v", err, expected)
	}
}

func mustParseURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	return parsed
}
