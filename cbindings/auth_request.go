package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"strings"

	"go.chrastecky.dev/bitsailor-core/bitwarden"
	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenFetchAuthRequest
func BitwardenFetchAuthRequest(
	client C.ClientHandle,
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	id C.UUID,
	requestOut *C.Handle,
	phraseOut **C.char,
) C.BitwardenResult {
	if requestOut == nil {
		setLastError(nullPointerError("requestOut"))
		return BitwardenError
	}
	if phraseOut == nil {
		setLastError(nullPointerError("phraseOut"))
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
	vaultGo, err := getHandleObj[bitwarden.Vault](handle(vault))
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

	fingerprint, err := crypto.FingerprintPhrase(vaultGo.GetVaultData().Profile.Email, authRequest.PublicKey)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	handleID := registerHandle(authRequest)
	*requestOut = C.Handle(handleID)
	*phraseOut = C.CString(strings.Join(fingerprint, "-"))

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
