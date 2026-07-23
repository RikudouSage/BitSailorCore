package main

/*
#include "bw_common.h"
*/
import "C"

//export BitwardenDeleteItem
func BitwardenDeleteItem(vault C.VaultHandle, ctx C.ContextHandle, session C.SessionHandle, itemID C.UUID) C.BitwardenResult {
	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	itemIDGo := parseUUIDFromC(itemID)
	err = vaultGo.DeleteItem(ctxGo, sessionGo, itemIDGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
