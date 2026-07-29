package bitwarden

import (
	"net/url"
	"testing"
)

func TestURLWithPathPreservesBasePath(t *testing.T) {
	t.Parallel()

	baseURL, err := url.Parse("https://example.com/bitwarden")
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	actual := urlWithPath(baseURL, "/identity/connect/token")
	expected := "https://example.com/bitwarden/identity/connect/token"
	if actual.String() != expected {
		t.Fatalf("urlWithPath() = %s, want %s", actual.String(), expected)
	}
	if baseURL.String() != "https://example.com/bitwarden" {
		t.Fatalf("base URL was mutated: %s", baseURL.String())
	}
}

func TestURLWithPathHandlesTrailingAndLeadingSlashes(t *testing.T) {
	t.Parallel()

	baseURL, err := url.Parse("https://example.com/bitwarden/")
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	actual := urlWithPath(baseURL, "ciphers")
	expected := "https://example.com/bitwarden/ciphers"
	if actual.String() != expected {
		t.Fatalf("urlWithPath() = %s, want %s", actual.String(), expected)
	}
}

func TestURLWithPathWithoutBasePath(t *testing.T) {
	t.Parallel()

	baseURL, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatalf("url.Parse() returned error: %v", err)
	}

	actual := urlWithPath(baseURL, "/api/sync")
	expected := "https://example.com/api/sync"
	if actual.String() != expected {
		t.Fatalf("urlWithPath() = %s, want %s", actual.String(), expected)
	}
}
