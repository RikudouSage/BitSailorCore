package bitwarden

type Vault interface {
}

type vault struct {
}

func newVault() *vault {
	return &vault{}
}
