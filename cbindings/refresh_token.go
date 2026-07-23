package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"go.chrastecky.dev/bitwarden-client/bitwarden"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

//export BitwardenRefreshToken
func BitwardenRefreshToken(client C.ClientHandle, ctx C.ContextHandle, session C.SessionHandle) C.BitwardenResult {
	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	ctxGo, err := getHandleObj[*contextHandle](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if err = clientGo.Auth().RefreshToken(ctxGo.ctx, sessionGo); err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
