package result

import (
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

type SendText struct {
	Text   string `json:"text"`
	Hidden bool   `json:"hidden"`
}

type SendFile struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	Size     string `json:"size"`
	SizeName string `json:"sizeName"`
}

type Send struct {
	ID             uuid.UUID      `json:"id"`
	AccessID       string         `json:"AccessId"`
	AuthType       SendAuthType   `json:"authType"`
	Name           string         `json:"name"`
	Disabled       bool           `json:"disabled"`
	RevisionDate   time.Time      `json:"revisionDate"`
	DeletionDate   time.Time      `json:"deletionDate"`
	HideEmail      bool           `json:"hideEmail"`
	Notes          *string        `json:"notes"`
	File           *SendFile      `json:"file"`
	Key            string         `json:"key"`
	AccessCount    uint           `json:"accessCount"`
	Password       *string        `json:"password"`
	ExpirationDate time.Time      `json:"expirationDate"`
	Type           SendType       `json:"type"`
	MaxAccessCount *uint          `json:"maxAccessCount"`
	Emails         types.CSVSlice `json:"emails"`
	Text           *SendText      `json:"text"`
}
