package bitwarden

import "go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"

type preLoginResponse struct {
	KDFType        crypto.KDFType `json:"kdf"`
	KDFIterations  int            `json:"kdfIterations"`
	KDFMemory      *int           `json:"kdfMemory"`
	KDFParallelism *int           `json:"kdfParallelism"`
}

type twoFactorErrorResponse struct {
	Error                   string   `json:"error"`
	ErrorDescription        string   `json:"error_description"`
	TwoFactorProviders      []string `json:"TwoFactorProviders"`
	SsoEmail2faSessionToken string   `json:"SsoEmail2faSessionToken"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`

	Key        string  `json:"Key"`
	PrivateKey *string `json:"PrivateKey"`

	KDFType        crypto.KDFType `json:"Kdf"`
	KDFIterations  int            `json:"KdfIterations"`
	KDFMemory      *int           `json:"KdfMemory"`
	KDFParallelism *int           `json:"KdfParallelism"`

	AccountKeys *struct {
		Object                     string `json:"Object"`
		PublicKeyEncryptionKeyPair *struct {
			Object            string  `json:"Object"`
			PublicKey         *string `json:"publicKey"`
			WrappedPrivateKey *string `json:"wrappedPrivateKey"`
		} `json:"publicKeyEncryptionKeyPair"`
	} `json:"AccountKeys"`
	UserDecryptionOptions *struct {
		Object               string `json:"Object"`
		HasMasterPassword    bool   `json:"HasMasterPassword"`
		MasterPasswordUnlock *struct {
			KDF *struct {
				Iterations  int            `json:"Iterations"`
				KDFType     crypto.KDFType `json:"KdfType"`
				Memory      *int           `json:"Memory"`
				Parallelism *int           `json:"Parallelism"`
			} `json:"Kdf"`
			MasterKeyEncryptedUserKey string `json:"MasterKeyEncryptedUserKey"`
			MasterKeyWrappedUserKey   string `json:"MasterKeyWrappedUserKey"`
			Salt                      string `json:"Salt"`
		} `json:"MasterPasswordUnlock"`
	} `json:"UserDecryptionOptions"`

	ForcePasswordReset   bool `json:"ForcePasswordReset"`
	MasterPasswordPolicy *struct {
		Object string `json:"Object"`
	} `json:"MasterPasswordPolicy"`
	ResetMasterPassword bool `json:"ResetMasterPassword"`
}

func (receiver *tokenResponse) GetPrivateKey() *string {
	if receiver.AccountKeys != nil && receiver.AccountKeys.PublicKeyEncryptionKeyPair != nil && receiver.AccountKeys.PublicKeyEncryptionKeyPair.WrappedPrivateKey != nil {
		return receiver.AccountKeys.PublicKeyEncryptionKeyPair.WrappedPrivateKey
	}

	return receiver.PrivateKey
}

func (receiver *tokenResponse) GetUserKey() *string {
	if receiver.UserDecryptionOptions != nil && receiver.UserDecryptionOptions.MasterPasswordUnlock != nil && receiver.UserDecryptionOptions.MasterPasswordUnlock.MasterKeyWrappedUserKey != "" {
		return new(receiver.UserDecryptionOptions.MasterPasswordUnlock.MasterKeyWrappedUserKey)
	}

	if receiver.Key == "" {
		return nil
	}

	return new(receiver.Key)
}
