package bitwarden

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewClientRequiresDeviceID(t *testing.T) {
	t.Parallel()

	_, err := NewClient()
	if err == nil {
		t.Fatal("NewClient() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "device ID is required") {
		t.Fatalf("NewClient() error = %v, want device ID validation error", err)
	}
}

func TestNewClientAppliesDefaultsAndOptions(t *testing.T) {
	t.Parallel()

	deviceID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	httpClient := &http.Client{}
	clientInterface, err := NewClient(
		WithDeviceID(deviceID),
		WithHTTPClient(httpClient),
		WithDebugLogsEnabled(true),
	)
	if err != nil {
		t.Fatalf("NewClient() returned error: %v", err)
	}

	client := clientInterface.(*client)
	if client.httpClient != httpClient {
		t.Fatal("client.httpClient was not set from option")
	}
	if !client.debugLogs {
		t.Fatal("client.debugLogs = false, want true")
	}
	if client.identityURL.String() != "https://vault.bitwarden.com" {
		t.Fatalf("identityURL = %s, want https://vault.bitwarden.com", client.identityURL)
	}
	if client.apiURL.String() != "https://api.bitwarden.com" {
		t.Fatalf("apiURL = %s, want https://api.bitwarden.com", client.apiURL)
	}
	if client.notificationsURL.String() != "https://notifications.bitwarden.com" {
		t.Fatalf("notificationsURL = %s, want https://notifications.bitwarden.com", client.notificationsURL)
	}
}

func TestWithBaseURLNormalizesSpecialAndSelfHostedURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		baseURL           string
		wantIdentity      string
		wantAPI           string
		wantNotifications string
	}{
		{
			name:              "US cloud",
			baseURL:           "https://bitwarden.com",
			wantIdentity:      "https://vault.bitwarden.com",
			wantAPI:           "https://api.bitwarden.com",
			wantNotifications: "https://notifications.bitwarden.com",
		},
		{
			name:              "EU cloud",
			baseURL:           "https://bitwarden.eu",
			wantIdentity:      "https://vault.bitwarden.eu",
			wantAPI:           "https://api.bitwarden.eu",
			wantNotifications: "https://notifications.bitwarden.eu",
		},
		{
			name:              "self hosted",
			baseURL:           "https://bw.example.test",
			wantIdentity:      "https://bw.example.test",
			wantAPI:           "https://bw.example.test/api",
			wantNotifications: "https://bw.example.test/notifications",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{}
			if err := WithBaseURL(test.baseURL)(client); err != nil {
				t.Fatalf("WithBaseURL() returned error: %v", err)
			}

			if client.identityURL.String() != test.wantIdentity {
				t.Fatalf("identityURL = %s, want %s", client.identityURL, test.wantIdentity)
			}
			if client.apiURL.String() != test.wantAPI {
				t.Fatalf("apiURL = %s, want %s", client.apiURL, test.wantAPI)
			}
			if client.notificationsURL.String() != test.wantNotifications {
				t.Fatalf("notificationsURL = %s, want %s", client.notificationsURL, test.wantNotifications)
			}
		})
	}
}

func TestURLSetterOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(*client) error
		want func(*client) string
	}{
		{
			name: "identity URL",
			run:  WithIdentityURL("https://identity.example.test"),
			want: func(client *client) string { return client.identityURL.String() },
		},
		{
			name: "API URL",
			run:  WithAPIURL("https://api.example.test"),
			want: func(client *client) string { return client.apiURL.String() },
		},
		{
			name: "notifications URL",
			run:  WithNotificationsURL("https://notifications.example.test"),
			want: func(client *client) string { return client.notificationsURL.String() },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{}
			if err := test.run(client); err != nil {
				t.Fatalf("option returned error: %v", err)
			}
			if got := test.want(client); !strings.HasPrefix(got, "https://") {
				t.Fatalf("configured URL = %s, want https URL", got)
			}
		})
	}
}

func TestURLSetterOptionsClearValues(t *testing.T) {
	t.Parallel()

	client := &client{}
	if err := WithBaseURL("https://bw.example.test")(client); err != nil {
		t.Fatalf("WithBaseURL() returned error: %v", err)
	}
	if err := WithIdentityURL("")(client); err != nil {
		t.Fatalf("WithIdentityURL(empty) returned error: %v", err)
	}
	if client.identityURL != nil {
		t.Fatalf("identityURL = %s, want nil", client.identityURL)
	}
	if err := WithAPIURL("")(client); err != nil {
		t.Fatalf("WithAPIURL(empty) returned error: %v", err)
	}
	if client.apiURL != nil {
		t.Fatalf("apiURL = %s, want nil", client.apiURL)
	}
	if err := WithNotificationsURL("")(client); err != nil {
		t.Fatalf("WithNotificationsURL(empty) returned error: %v", err)
	}
	if client.notificationsURL != nil {
		t.Fatalf("notificationsURL = %s, want nil", client.notificationsURL)
	}
}

func TestClientAccessorsCacheServices(t *testing.T) {
	t.Parallel()

	deviceID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	clientInterface, err := NewClient(WithDeviceID(deviceID))
	if err != nil {
		t.Fatalf("NewClient() returned error: %v", err)
	}
	client := clientInterface.(*client)

	if client.Auth() != client.Auth() {
		t.Fatal("Auth() did not cache service")
	}
	if client.Vault() != client.Vault() {
		t.Fatal("Vault() did not cache service")
	}
	if client.Notifications() != client.Notifications() {
		t.Fatal("Notifications() did not cache service")
	}
}
