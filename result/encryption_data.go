package result

import (
	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/dto"
)

type AccountKeys struct {
	PublicKey         *string
	WrappedPrivateKey *string
}

type EncryptionData struct {
	UserKey          dto.Key               `json:"userKey"`
	OrganizationKeys map[uuid.UUID]dto.Key `json:"organizationKeys"`

	EncryptedUserKey    *string `json:"encryptedUserKey"`
	EncryptedPrivateKey *string `json:"encryptedPrivateKey"`

	KDFType        crypto.KDFType `json:"kdfType"`
	KDFIterations  int            `json:"kdfIterations"`
	KDFMemory      *int           `json:"kdfMemory"`
	KDFParallelism *int           `json:"kdfParallelism"`
}
