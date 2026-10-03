package main

/*
#include "bw_common.h"
#include "bw_auth.h"
*/
import "C"

import "go.chrastecky.dev/bitsailor-core/bitwarden/dto"

//export BitwardenInitializeTOTPProvider
func BitwardenInitializeTOTPProvider(
	client C.ClientHandle,
	ctx C.ContextHandle,
	email, password *C.char,
	kind C.BitwardenTfaKind,
) C.BitwardenResult {
	clientGo, ctxGo, err := getCommonAuthHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	err = clientGo.Auth().InitializeTOTPProvider(
		ctxGo,
		C.GoString(email),
		C.GoString(password),
		dto.TFAKind(kind),
	)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
