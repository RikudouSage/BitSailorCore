package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"errors"
	"fmt"

	"go.chrastecky.dev/bitsailor-core/bitwarden"
)

var ErrNoVault = errors.New("no vault")
var ErrNoVaultData = errors.New("no vault data")

//export BitwardenGetEmail
func BitwardenGetEmail(vault C.VaultHandle, out **C.char) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	vaultGo, err := getHandleObj[bitwarden.Vault](handle(vault))
	if err != nil {
		setLastError(fmt.Errorf("%w: %w", ErrNoVault, err))
		return BitwardenError
	}

	vaultData := vaultGo.GetVaultData()
	if vaultData == nil || vaultData.Profile == nil {
		setLastError(ErrNoVaultData)
		return BitwardenError
	}

	*out = C.CString(vaultData.Profile.Email)
	clearLastError()
	return BitwardenSuccess
}
