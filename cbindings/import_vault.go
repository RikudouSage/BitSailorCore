package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"encoding/json"

	"go.chrastecky.dev/bitsailor-core/bitwarden"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenImportVault
func BitwardenImportVault(inVault C.VaultHandle, exportData *C.char, outVault *C.VaultHandle) C.BitwardenResult {
	if outVault == nil {
		setLastError(nullPointerError("outVault"))
		return BitwardenError
	}

	vault, err := getHandleObj[bitwarden.Vault](handle(inVault))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	var data *result.VaultData
	if err = json.Unmarshal([]byte(C.GoString(exportData)), &data); err != nil {
		setLastError(err)
		return BitwardenError
	}

	newVault := vault.WithVaultData(data)
	handleID := registerHandle(newVault)

	*outVault = C.VaultHandle(handleID)

	clearLastError()
	return BitwardenSuccess
}
