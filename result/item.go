package result

import (
	"time"

	"github.com/google/uuid"
	clone "github.com/huandu/go-clone/generic"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/types"
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

type URIMatchType int

const (
	URIMatchTypeDomain URIMatchType = iota
	URIMatchTypeHost
	URIMatchTypeStartsWith
	URIMatchTypeExact
	URIMatchTypeRegularExpression
	URIMatchTypeNever
)

type FieldType int

const (
	FieldTypeText FieldType = iota
	FieldTypeHidden
	FieldTypeCheckbox
	FieldTypeLinkedID
)

type ItemPermissions struct {
	Delete  bool `json:"delete"`
	Restore bool `json:"restore"`
}

type ItemLoginURI struct {
	URI         string       `json:"uri"`
	URIChecksum string       `json:"uriChecksum"`
	Match       URIMatchType `json:"match"` // todo
}

type ItemLogin struct {
	URI                  string          `json:"uri"`
	URIs                 []*ItemLoginURI `json:"uris"`
	Username             *string         `json:"username"`
	Password             *string         `json:"password"`
	PasswordRevisionDate *time.Time      `json:"passwordRevisionDate"`
	TOTP                 *string         `json:"totp"`
	AutofillOnPageLoad   any             `json:"autofillOnPageLoad"` // todo
	Fido2Credentials     any             `json:"fido2Credentials"`   // todo
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

type Field struct {
	Type     FieldType `json:"type"`
	Name     string    `json:"name"`
	Value    *string   `json:"value"`
	LinkedID *int      `json:"linkedId"`
}

func (receiver Field) CheckboxValue() bool {
	return receiver.Type == FieldTypeCheckbox && receiver.Value != nil && *receiver.Value == "true"
}

type Item struct {
	ID                  uuid.UUID        `json:"id,omitzero"`
	Type                ItemType         `json:"type"`
	Notes               *string          `json:"notes,omitempty"`
	OrganizationUseTOTP *bool            `json:"organizationUseTotp,omitempty"`
	RevisionDate        time.Time        `json:"revisionDate,omitzero"`
	DeletedDate         *time.Time       `json:"deletedDate,omitempty"`
	Favorite            bool             `json:"favorite"`
	OrganizationID      uuid.UUID        `json:"organizationId,omitzero"`
	Key                 *string          `json:"key,omitempty"`
	Permissions         *ItemPermissions `json:"permissions,omitempty"`
	Edit                bool             `json:"edit"`
	CollectionIDs       []uuid.UUID      `json:"collectionIds,omitempty"`
	ArchivedDate        *time.Time       `json:"archivedDate,omitempty"`
	FolderID            uuid.UUID        `json:"folderId,omitzero"`
	ViewPassword        bool             `json:"viewPassword"`
	Name                string           `json:"name"`
	CreationDate        time.Time        `json:"creationDate,omitzero"`
	Reprompt            types.NumBool    `json:"reprompt"`
	Fields              []*Field         `json:"fields,omitempty"`

	Login      *ItemLogin      `json:"login,omitempty"`
	Card       *ItemCard       `json:"card,omitempty"`
	SecureNote *ItemSecureNote `json:"secureNote,omitempty"`
	Identity   *ItemIdentity   `json:"identity,omitempty"`
	SSHKey     *ItemSSHKey     `json:"sshKey,omitempty"`

	DecryptionError error

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

func (receiver *Item) AsInvalidItem(err error) *Item {
	fields := make([]*Field, 0, len(receiver.Fields))
	for _, field := range receiver.Fields {
		if field == nil {
			fields = append(fields, nil)
			continue
		}

		fields = append(fields, &Field{
			Type:     field.Type,
			LinkedID: field.LinkedID,
		})
	}

	var secureNote *ItemSecureNote
	if receiver.SecureNote != nil {
		secureNote = &ItemSecureNote{Type: receiver.SecureNote.Type}
	}

	return &Item{
		ID:                  receiver.ID,
		Type:                receiver.Type,
		OrganizationUseTOTP: receiver.OrganizationUseTOTP,
		RevisionDate:        receiver.RevisionDate,
		DeletedDate:         receiver.DeletedDate,
		Favorite:            receiver.Favorite,
		OrganizationID:      receiver.OrganizationID,
		Permissions:         receiver.Permissions,
		Edit:                receiver.Edit,
		CollectionIDs:       clone.Clone(receiver.CollectionIDs),
		ArchivedDate:        receiver.ArchivedDate,
		FolderID:            receiver.FolderID,
		ViewPassword:        receiver.ViewPassword,
		CreationDate:        receiver.CreationDate,
		Reprompt:            receiver.Reprompt,
		Fields:              fields,
		SecureNote:          secureNote,
		DecryptionError:     err,
	}
}
