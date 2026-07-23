package bitwarden

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

var ErrMissingVault = errors.New("the vault is missing, either run Sync() or WithVaultData() on it")
var ErrItemNotFound = errors.New("the item was not found in the vault")

type Vault interface {
	Sync(ctx context.Context, session *result.Session) (Vault, error)
	GetItem(ctx context.Context, session *result.Session, itemID uuid.UUID) (*result.Item, error)
	GetItems(ctx context.Context, session *result.Session) ([]*result.Item, error)
	CreateItem(ctx context.Context, session *result.Session, item *result.Item) error
	DeleteItem(ctx context.Context, session *result.Session, itemID uuid.UUID) error

	GetVaultData() *result.VaultData
	WithVaultData(vaultData *result.VaultData) Vault
}

type vault struct {
	baseURL    *url.URL
	httpClient *http.Client

	vaultData *result.VaultData
}

func newVault(baseURL *url.URL, httpClient *http.Client) *vault {
	return &vault{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (receiver *vault) GetVaultData() *result.VaultData {
	return receiver.vaultData
}

func (receiver *vault) WithVaultData(vaultData *result.VaultData) Vault {
	clone := new(*receiver)
	clone.vaultData = vaultData

	return clone
}
