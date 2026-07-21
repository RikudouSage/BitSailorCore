package result

import "github.com/google/uuid"

type CollectionType int // todo

type Collection struct {
	ID                         uuid.UUID      `json:"id"`
	OrganizationID             uuid.UUID      `json:"organizationId"`
	Type                       CollectionType `json:"type"`
	DefaultUserCollectionEmail *string        `json:"defaultUserCollectionEmail"`
	ReadOnly                   bool           `json:"readOnly"`
	HidePasswords              bool           `json:"hidePasswords"`
	Manage                     bool           `json:"manage"`
	Name                       string         `json:"name"`
	ExternalID                 *string        `json:"externalId"`
}
