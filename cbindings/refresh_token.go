package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenRefreshToken
func BitwardenRefreshToken(client C.ClientHandle, ctx C.ContextHandle, session C.SessionHandle) C.BitwardenResult {
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

	if err = clientGo.Auth().RefreshToken(ctxGo, sessionGo); err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
