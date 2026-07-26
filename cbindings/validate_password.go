package main

/*
#include "bw_common.h"
*/
import "C"

import (
	"go.chrastecky.dev/bitsailor-core/bitwarden"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenValidatePassword
func BitwardenValidatePassword(
	client C.ClientHandle,
	ctx C.ContextHandle,
	email, password *C.char,
	session C.SessionHandle,
) C.BitwardenResult {
	if email == nil {
		setLastError(nullPointerError("email"))
		return BitwardenError
	}
	if password == nil {
		setLastError(nullPointerError("password"))
		return BitwardenError
	}

	ctxGo, err := getHandleObj[*contextHandle](handle(ctx))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	sessionClone := &result.Session{
		Auth: sessionGo.Auth,
		Encryption: &result.EncryptionData{
			UserKey:             nil,
			OrganizationKeys:    sessionGo.Encryption.OrganizationKeys,
			EncryptedUserKey:    sessionGo.Encryption.EncryptedUserKey,
			EncryptedPrivateKey: sessionGo.Encryption.EncryptedPrivateKey,
			KDFType:             sessionGo.Encryption.KDFType,
			KDFIterations:       sessionGo.Encryption.KDFIterations,
			KDFMemory:           sessionGo.Encryption.KDFMemory,
			KDFParallelism:      sessionGo.Encryption.KDFParallelism,
		},
	}

	if err = clientGo.Auth().UnlockSession(ctxGo.ctx, sessionClone, C.GoString(email), C.GoString(password)); err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
