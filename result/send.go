package result

import (
	"io"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/types"
)

type SendAuthType uint8
type SendType uint8

const (
	SendAuthTypeSpecificPeople SendAuthType = iota
	SendAuthTypePassword
	SendAuthTypeNoAuth
)

const (
	SendTypeText SendType = iota
	SendTypeFile
)

type SendText struct {
	Text   string `json:"text"`
	Hidden bool   `json:"hidden"`
}

type SendFile struct {
	ID       string `json:"id,omitempty"`
	FileName string `json:"fileName"`
	Size     string `json:"size,omitempty"`
	SizeName string `json:"sizeName,omitempty"`
}

type Send struct {
	ID             uuid.UUID      `json:"id,omitzero"`
	AccessID       *string        `json:"accessId,omitzero"`
	AuthType       SendAuthType   `json:"authType"`
	Name           string         `json:"name"`
	Disabled       bool           `json:"disabled"`
	RevisionDate   time.Time      `json:"revisionDate,omitzero"`
	DeletionDate   time.Time      `json:"deletionDate"`
	HideEmail      bool           `json:"hideEmail"`
	Notes          *string        `json:"notes,omitempty"`
	File           *SendFile      `json:"file,omitempty"`
	Key            string         `json:"key,omitzero"`
	AccessCount    uint           `json:"accessCount,omitzero"`
	Password       *string        `json:"password,omitempty"`
	ExpirationDate time.Time      `json:"expirationDate,omitzero"`
	Type           SendType       `json:"type"`
	MaxAccessCount *uint          `json:"maxAccessCount,omitempty"`
	Emails         types.CSVSlice `json:"emails,omitempty"`
	Text           *SendText      `json:"text,omitempty"`

	// only for creation
	FileLength int       `json:"fileLength,omitzero"`
	InputFile  io.Reader `json:"-"`
}
