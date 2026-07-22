package result

import (
	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/dto"
)

type AccountKeys struct {
	PublicKey         *string
	WrappedPrivateKey *string
}

type EncryptionData struct {
	UserKey          dto.Key
	OrganizationKeys map[uuid.UUID]dto.Key

	EncryptedUserKey    *string
	EncryptedPrivateKey *string

	KDFType        crypto.KDFType
	KDFIterations  int
	KDFMemory      *int
	KDFParallelism *int
}
