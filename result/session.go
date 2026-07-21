package result

import (
	"errors"
)

type Session struct {
	Auth       *AuthData
	Encryption *EncryptionData
}

func (receiver *Session) ValidateForUnlock() error {
	if receiver == nil {
		return errors.New("locked session is nil")
	}
	if receiver.Auth == nil {
		return errors.New("locked session auth data is nil")
	}
	if receiver.Encryption.EncryptedUserKey == nil {
		return errors.New("locked session missing encrypted user key")
	}

	return nil
}
