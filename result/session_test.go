package result

import (
	"strings"
	"testing"
)

func TestSessionValidateForUnlock(t *testing.T) {
	t.Parallel()

	encryptedUserKey := "2.iv|ciphertext|mac"
	tests := []struct {
		name    string
		session *Session
		wantErr string
	}{
		{
			name:    "nil session",
			session: nil,
			wantErr: "locked session is nil",
		},
		{
			name:    "nil auth",
			session: &Session{},
			wantErr: "locked session auth data is nil",
		},
		{
			name:    "nil encryption",
			session: &Session{Auth: &AuthData{}},
			wantErr: "locked session encryption data is nil",
		},
		{
			name:    "missing encrypted user key",
			session: &Session{Auth: &AuthData{}, Encryption: &EncryptionData{}},
			wantErr: "locked session missing encrypted user key",
		},
		{
			name:    "valid",
			session: &Session{Auth: &AuthData{}, Encryption: &EncryptionData{EncryptedUserKey: &encryptedUserKey}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.session.ValidateForUnlock()
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("ValidateForUnlock() returned error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("ValidateForUnlock() error = %v, want %q", err, test.wantErr)
			}
		})
	}
}
