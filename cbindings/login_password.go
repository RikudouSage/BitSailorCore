package main

/*
#include <bw_common.h>
*/
import "C"
import "go.chrastecky.dev/bitwarden-client/bitwarden"

//export BitwardenLoginPassword
func BitwardenLoginPassword(
	client C.ClientHandle,
	ctx C.ContextHandle,
	email, password *C.char,
	outHandle *C.SessionHandle,
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

	emailStr := C.GoString(email)
	passwordStr := C.GoString(password)

	session, err := clientGo.Auth().LoginPassword(ctxGo.ctx, emailStr, passwordStr)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sessionHandleID := registerHandle(session)
	*outHandle = C.SessionHandle(sessionHandleID)

	clearLastError()
	return BitwardenSuccess
}
