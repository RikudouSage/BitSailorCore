package main

/*
#include "bw_common.h"
#include "bw_item.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenGetItems
func BitwardenGetItems(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	out *C.BitwardenItemSlice,
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

	items, err := vaultGo.GetItems(ctxGo, sessionGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*out = bitwardenItemSliceIntoC(items)

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenFreeItems
func BitwardenFreeItems(items *C.BitwardenItemSlice) {
	freeBitwardenItemSlice(items)
}

func bitwardenItemSliceIntoC(items []*result.Item) C.BitwardenItemSlice {
	if len(items) == 0 {
		return C.BitwardenItemSlice{}
	}

	cItems := (*C.BitwardenItem)(C.malloc(C.size_t(len(items)) * C.size_t(unsafe.Sizeof(C.BitwardenItem{}))))
	out := unsafe.Slice(cItems, len(items))
	for i, item := range items {
		out[i] = bitwardenItemIntoC(item)
	}

	return C.BitwardenItemSlice{items: cItems, len: C.size_t(len(items))}
}

func freeBitwardenItemSlice(items *C.BitwardenItemSlice) {
	if items == nil {
		return
	}

	cItems := unsafe.Slice(items.items, int(items.len))
	for i := range cItems {
		freeBitwardenItem(&cItems[i])
	}
	C.free(unsafe.Pointer(items.items))

	*items = C.BitwardenItemSlice{}
}
