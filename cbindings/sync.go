package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenSyncVault
func BitwardenSyncVault(
	client C.ClientHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	outHandle *C.VaultHandle,
) C.BitwardenResult {
	if outHandle == nil {
		setLastError(nullPointerError("outHandle"))
		return BitwardenError
	}

	clientGo, ctxGo, err := getCommonAuthHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	syncResult, err := clientGo.Vault().Sync(ctxGo, sessionGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	outHandleGo := registerHandle(syncResult)
	*outHandle = C.VaultHandle(outHandleGo)

	clearLastError()
	return BitwardenSuccess
}
