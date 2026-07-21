package result

import (
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/types"
)

type ItemType int // todo

const (
	ItemTypeLogin ItemType = iota + 1
	ItemTypeSecureNote
	ItemTypeCard
	ItemTypIdentity
	ItemTypeSSHKey
	ItemTypeBankAccount
	ItemTypeDriversLicense
	ItemTypePassport
)

type ItemPermissions struct {
	Delete  bool `json:"delete"`
	Restore bool `json:"restore"`
}

type ItemLogin struct {
	URI  string `json:"uri"`
	URIs []struct {
		URI         string `json:"uri"`
		URIChecksum string `json:"uriChecksum"`
		Match       any    `json:"match"` // todo
	} `json:"uris"`
	Username             string    `json:"username"`
	Password             string    `json:"password"`
	PasswordRevisionDate time.Time `json:"passwordRevisionDate"`
	TOTP                 *string   `json:"totp"`
	AutofillOnPageLoad   any       `json:"autofillOnPageLoad"` // todo
	Fido2Credentials     any       `json:"fido2Credentials"`   // todo
}

type ItemCard struct {
	CardholderName  string `json:"cardholderName"`
	Brand           string `json:"brand"`
	Number          string `json:"number"`
	ExpirationMonth string `json:"expMonth"`
	ExpirationYear  string `json:"expYear"`
	Code            string `json:"code"`
}

type ItemSecureNote struct {
	Type int `json:"type"` // todo what is this
}

type ItemIdentity struct {
	FirstName      *string `json:"firstName"`
	MiddleName     *string `json:"middleName"`
	LastName       *string `json:"lastName"`
	Title          *string `json:"title"`
	PassportNumber *string `json:"passportNumber"`

	Username *string `json:"username"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`

	AddressLine1 *string `json:"addressLine1"`
	AddressLine2 *string `json:"addressLine2"`
	AddressLine3 *string `json:"addressLine3"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	PostalCode   *string `json:"postalCode"`
	Country      *string `json:"country"`
	SSN          *string `json:"ssn"`

	Company *string `json:"company"`
}

type ItemSSHKey struct {
	PrivateKey     string `json:"privateKey"`
	PublicKey      string `json:"publicKey"`
	KeyFingerprint string `json:"keyFingerprint"`
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

	Login      *ItemLogin      `json:"login"`
	Card       *ItemCard       `json:"card"`
	SecureNote *ItemSecureNote `json:"secureNote"`
	Identity   *ItemIdentity   `json:"identity"`
	SSHKey     *ItemSSHKey     `json:"sshKey"`

	// bankAccount
	// identity
	// passport
	// attachments
	// data ?
	// passwordHistory
	// driversLicense
	// sshKey
	// fields
}
