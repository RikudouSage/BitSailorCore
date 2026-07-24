package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenUnlockSession
func BitwardenUnlockSession(
	client C.ClientHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	email, password *C.char,
) C.BitwardenResult {
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
	emailGo := C.GoString(email)
	passwordGo := C.GoString(password)

	err = clientGo.Auth().UnlockSession(ctxGo, sessionGo, emailGo, passwordGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
