package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"context"

	"go.chrastecky.dev/bitsailor-core/bitwarden"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func getCommonVaultHandles(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
) (vaultGo bitwarden.Vault, ctxGo context.Context, sessionGo *result.Session, err error) {
	vaultGo, err = getHandleObj[bitwarden.Vault](handle(vault))
	if err != nil {
		return
	}
	ctxWrapper, err := getHandleObj[*contextHandle](handle(ctx))
	if err != nil {
		return
	}
	ctxGo = ctxWrapper.ctx

	sessionGo, err = getHandleObj[*result.Session](handle(session))
	if err != nil {
		return
	}

	return
}

func getCommonAuthHandles(
	client C.ClientHandle,
	ctx C.ContextHandle,
) (clientGo bitwarden.Client, ctxGo context.Context, err error) {
	clientGo, err = getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		return
	}
	ctxWrapper, err := getHandleObj[*contextHandle](handle(ctx))
	if err != nil {
		return
	}
	ctxGo = ctxWrapper.ctx

	return
}
