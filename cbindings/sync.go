package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"go.chrastecky.dev/bitwarden-client/bitwarden"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

//export BitwardenSyncVault
func BitwardenSyncVault(
	client C.ClientHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	outHandle *C.VaultHandle,
) C.Result {
	if outHandle == nil {
		setLastError(nullPointerError("outHandle"))
		return BitwardenError
	}

	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	ctxGo, err := getHandleObj[*contextHandle](handle(ctx))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	syncResult, err := clientGo.Vault().Sync(ctxGo.ctx, sessionGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	outHandleGo := registerHandle(syncResult)
	*outHandle = C.VaultHandle(outHandleGo)

	clearLastError()
	return BitwardenSuccess
}
