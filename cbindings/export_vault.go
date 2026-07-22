package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"encoding/json"

	"go.chrastecky.dev/bitwarden-client/bitwarden"
)

//export BitwardenExportEncryptedVault
func BitwardenExportEncryptedVault(vault C.VaultHandle, jsonOut **C.char) C.Result {
	if jsonOut == nil {
		setLastError(nullPointerError("jsonOut"))
		return BitwardenError
	}

	vaultGo, err := getHandleObj[bitwarden.Vault](handle(vault))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	outBytes, err := json.Marshal(vaultGo.GetVaultData())
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*jsonOut = C.CString(string(outBytes))

	clearLastError()
	return BitwardenSuccess
}
