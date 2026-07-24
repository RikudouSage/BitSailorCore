package main

/*
#include "bw_common.h"
*/
import "C"
import "go.chrastecky.dev/bitsailor-core/bitwarden"

//export BitwardenGetVault
func BitwardenGetVault(client C.ClientHandle, out *C.VaultHandle) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	handleID := registerHandle(clientGo.Vault())
	*out = C.VaultHandle(handleID)

	clearLastError()
	return BitwardenSuccess
}
