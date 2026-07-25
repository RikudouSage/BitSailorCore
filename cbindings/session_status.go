package main

/*
#include "bw_common.h"

typedef enum {
	BitwardenSessionStatusUnlocked,
	BitwardenSessionStatusLocked,
	BitwardenSessionStatusNone,
} BitwardenSessionStatus;
*/
import "C"

import "go.chrastecky.dev/bitsailor-core/bitwarden/result"

//export BitwardenGetSessionStatus
func BitwardenGetSessionStatus(session C.SessionHandle, out *C.BitwardenSessionStatus) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if sessionGo.Auth == nil || sessionGo.Auth.AccessToken == "" {
		*out = C.BitwardenSessionStatusNone
		clearLastError()
		return BitwardenSuccess
	}

	if sessionGo.Encryption.UserKey != nil {
		*out = C.BitwardenSessionStatusUnlocked
		clearLastError()
		return BitwardenSuccess
	}

	*out = C.BitwardenSessionStatusLocked
	clearLastError()
	return BitwardenSuccess
}
