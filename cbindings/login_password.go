package main

/*
#include "bw_common.h"
*/
import "C"

//export BitwardenLoginPassword
func BitwardenLoginPassword(
	client C.ClientHandle,
	ctx C.ContextHandle,
	email, password *C.char,
	outHandle *C.SessionHandle,
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

	emailStr := C.GoString(email)
	passwordStr := C.GoString(password)

	session, err := clientGo.Auth().LoginPassword(ctxGo, emailStr, passwordStr)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sessionHandleID := registerHandle(session)
	*outHandle = C.SessionHandle(sessionHandleID)

	clearLastError()
	return BitwardenSuccess
}
