package main

/*
#include "bw_common.h"
#include "bw_send.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenGetSends
func BitwardenGetSends(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	out *C.BitwardenSendSlice,
) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sends, err := vaultGo.GetSends(ctxGo, sessionGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*out = bitwardenSendSliceIntoC(sends)

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenFreeSends
func BitwardenFreeSends(sends *C.BitwardenSendSlice) {
	freeBitwardenSendSlice(sends)
}

func bitwardenSendSliceIntoC(sends []*result.Send) C.BitwardenSendSlice {
	if len(sends) == 0 {
		return C.BitwardenSendSlice{}
	}

	cSends := (*C.BitwardenSend)(C.malloc(C.size_t(len(sends)) * C.size_t(unsafe.Sizeof(C.BitwardenSend{}))))
	out := unsafe.Slice(cSends, len(sends))
	for i, send := range sends {
		out[i] = bitwardenSendIntoC(send)
	}

	return C.BitwardenSendSlice{items: cSends, len: C.size_t(len(sends))}
}

func freeBitwardenSendSlice(sends *C.BitwardenSendSlice) {
	if sends == nil {
		return
	}

	cSends := unsafe.Slice(sends.items, int(sends.len))
	for i := range cSends {
		freeBitwardenSend(&cSends[i])
	}
	C.free(unsafe.Pointer(sends.items))

	*sends = C.BitwardenSendSlice{}
}
