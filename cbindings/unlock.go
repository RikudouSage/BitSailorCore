package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"encoding/base64"
	"errors"
	"fmt"

	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenUnlockSession
func BitwardenUnlockSession(
	client C.ClientHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	email, password *C.char,
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
	emailGo := C.GoString(email)
	passwordGo := C.GoString(password)

	err = clientGo.Auth().UnlockSession(ctxGo, sessionGo, emailGo, passwordGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenUnlockSessionWithUserKey
func BitwardenUnlockSessionWithUserKey(session C.SessionHandle, base64UserKey *C.char) C.BitwardenResult {
	if base64UserKey == nil {
		setLastError(nullPointerError("base64UserKey"))
		return BitwardenError
	}

	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if sessionGo.Encryption == nil {
		setLastError(errors.New("the session does not have any encryption data, please make sure to unlock an already populated session"))
		return BitwardenError
	}

	userKey, err := base64.StdEncoding.DecodeString(C.GoString(base64UserKey))
	if err != nil {
		setLastError(fmt.Errorf("failed base64 decoding the user key: %w", err))
		return BitwardenError
	}
	if len(userKey) != 64 {
		setLastError(fmt.Errorf("the session expects a 64 byte key, got %d instead", len(userKey)))
		return BitwardenError
	}
	sessionGo.Encryption.UserKey = userKey

	clearLastError()
	return BitwardenSuccess
}
