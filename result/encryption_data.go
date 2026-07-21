package result

import "go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"

type AccountKeys struct {
	PublicKey         *string
	WrappedPrivateKey *string
}

type EncryptionData struct {
	UserKey []byte

	EncryptedUserKey    *string
	EncryptedPrivateKey *string

	KDFType        crypto.KDFType
	KDFIterations  int
	KDFMemory      *int
	KDFParallelism *int
}
