package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"go.chrastecky.dev/bitwarden-client/bitwarden"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

//export BitwardenUnlockSession
func BitwardenUnlockSession(
	client C.ClientHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	email, password *C.char,
) C.Result {
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
	emailGo := C.GoString(email)
	passwordGo := C.GoString(password)

	err = clientGo.Auth().UnlockSession(ctxGo.ctx, sessionGo, emailGo, passwordGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
