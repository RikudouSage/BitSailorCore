package main

/*
#include "bw_common.h"
*/
import "C"
import "go.chrastecky.dev/bitwarden-client/bitwarden"

//export BitwardenLoginApiKey
func BitwardenLoginApiKey(
	client C.ClientHandle,
	ctx C.ContextHandle,
	clientID, clientSecret *C.char,
	outHandle *C.SessionHandle,
) C.BitwardenResult {
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

	clientIDGo := C.GoString(clientID)
	clientSecretGo := C.GoString(clientSecret)

	session, err := clientGo.Auth().LoginApiKey(ctxGo.ctx, clientIDGo, clientSecretGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sessionHandleID := registerHandle(session)
	*outHandle = C.SessionHandle(sessionHandleID)

	clearLastError()
	return BitwardenSuccess
}
