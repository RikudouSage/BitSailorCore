package main

/*
#include "bw_common.h"
#include "bw_auth.h"
*/
import "C"

//export BitwardenLoginPassword
func BitwardenLoginPassword(
	client C.ClientHandle,
	ctx C.ContextHandle,
	email, password *C.char,
	tfa *C.BitwardenTfaConfig,
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
	twoFa := goTFAConfigFromC(tfa)

	session, err := clientGo.Auth().LoginPassword(ctxGo, emailStr, passwordStr, twoFa)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sessionHandleID := registerHandle(session)
	*outHandle = C.SessionHandle(sessionHandleID)

	clearLastError()
	return BitwardenSuccess
}
