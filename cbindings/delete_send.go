package main

/*
#include "bw_common.h"
*/
import "C"

//export BitwardenDeleteSend
func BitwardenDeleteSend(vault C.VaultHandle, ctx C.ContextHandle, session C.SessionHandle, sendID C.UUID) C.BitwardenResult {
	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sendIDGo := parseUUIDFromC(sendID)
	err = vaultGo.DeleteSend(ctxGo, sessionGo, sendIDGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
