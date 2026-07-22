package result

import (
	"github.com/google/uuid"
)

type VaultData struct {
	Profile     *Profile      `json:"profile"`
	Folders     []*Folder     `json:"folders"`
	Collections []*Collection `json:"collections"`
	Items       []*Item       `json:"ciphers"`

	// policies
	// sends
	// domains
	// policiesNew
	// userDecryption
}

func (receiver VaultData) GetOrganizationKeys() map[uuid.UUID]string {
	result := make(map[uuid.UUID]string)
	for _, org := range receiver.Profile.Organizations {
		result[org.ID] = org.Key
	}

	return result
}
