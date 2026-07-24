package main

/*
#include "bw_common.h"
#include "bw_item.h"
*/
import "C"

//export BitwardenUpdateItem
func BitwardenUpdateItem(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	item *C.BitwardenItem,
	outItem *C.BitwardenItem,
) C.BitwardenResult {
	if item == nil {
		setLastError(nullPointerError("item"))
		return BitwardenError
	}
	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	itemGo := bitwardenItemFromC(item)
	err = vaultGo.UpdateItem(ctxGo, sessionGo, itemGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if outItem != nil {
		*outItem = bitwardenItemIntoC(itemGo)
	}

	clearLastError()
	return BitwardenSuccess
}
