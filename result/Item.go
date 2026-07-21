package result

import (
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/types"
)

type ItemType int // todo

type ItemPermissions struct {
	Delete  bool `json:"delete"`
	Restore bool `json:"restore"`
}

type ItemLogin struct {
}

type Item struct {
	ID                  uuid.UUID        `json:"id"`
	Type                ItemType         `json:"type"`
	Notes               *string          `json:"notes"`
	OrganizationUseTOTP *bool            `json:"organizationUseTotp"`
	RevisionDate        time.Time        `json:"revisionDate"`
	DeletedDate         *time.Time       `json:"deletedDate"`
	Favorite            bool             `json:"favorite"`
	OrganizationID      uuid.UUID        `json:"organizationId"`
	Key                 string           `json:"key"`
	Permissions         *ItemPermissions `json:"permissions"`
	Edit                bool             `json:"edit"`
	CollectionIDs       []uuid.UUID      `json:"collectionIds"`
	ArchivedDate        *time.Time       `json:"archivedDate"`
	FolderID            uuid.UUID        `json:"folderId"`
	ViewPassword        bool             `json:"viewPassword"`
	Name                string           `json:"name"`
	CreationDate        time.Time        `json:"creationDate"`
	Reprompt            types.NumBool    `json:"reprompt"`

	Login *ItemLogin `json:"login"`

	// bankAccount
	// notes
	// secureNote
	// identity
	// passport
	// attachments
	// card
	// data ?
	// passwordHistory
	// driversLicense
	// sshKey
	// fields
}
