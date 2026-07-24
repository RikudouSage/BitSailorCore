package main

/*
#include "bw_common.h"
#include "bw_send.h"
#include <stdlib.h>
*/
import "C"

//export BitwardenGetSend
func BitwardenGetSend(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	sendID C.UUID,
	outSend *C.BitwardenSend,
) C.BitwardenResult {
	if outSend == nil {
		setLastError(nullPointerError("outSend"))
		return BitwardenError
	}

	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sendIDGo := parseUUIDFromC(sendID)
	send, err := vaultGo.GetSend(ctxGo, sessionGo, sendIDGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*outSend = bitwardenSendIntoC(send)

	clearLastError()
	return BitwardenSuccess
}
