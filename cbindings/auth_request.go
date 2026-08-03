package main

/*
#include "bw_common.h"
*/
import "C"
import "go.chrastecky.dev/bitsailor-core/bitwarden/result"

//export BitwardenFetchAuthRequest
func BitwardenFetchAuthRequest(client C.ClientHandle, ctx C.ContextHandle, session C.SessionHandle, id C.UUID, out *C.Handle) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

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

	idGo := parseUUIDFromC(id)

	authRequest, err := clientGo.Auth().FetchAuthRequest(ctxGo, sessionGo, idGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	handleID := registerHandle(authRequest)
	*out = C.Handle(handleID)
	clearLastError()
	return BitwardenSuccess
}

//export BitwardenRespondToAuthRequest
func BitwardenRespondToAuthRequest(
	client C.ClientHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	authRequest C.Handle,
	approved C.bool,
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

	authRequestGo, err := getHandleObj[*result.AuthRequest](handle(authRequest))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	err = clientGo.Auth().RespondToAuthRequest(ctxGo, sessionGo, authRequestGo, bool(approved))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
