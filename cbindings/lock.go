package main

/*
#include "bw_common.h"
*/
import "C"
import "go.chrastecky.dev/bitsailor-core/bitwarden/result"

//export BitwardenLockSession
func BitwardenLockSession(session C.SessionHandle) C.BitwardenResult {
	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sessionGo.Encryption.UserKey = nil
	clearLastError()
	return BitwardenSuccess
}
