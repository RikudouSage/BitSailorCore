package main

/*
#include "bw_common.h"
#include "bw_item.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/samber/lo"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
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

	converted := lo.Map(items, func(item *result.Item, _ int) C.BitwardenItem {
		return bitwardenItemIntoC(item)
	})

	cItems := (*C.BitwardenItem)(C.malloc(C.size_t(len(converted)) * C.size_t(unsafe.Sizeof(C.BitwardenItem{}))))
	copy(unsafe.Slice(cItems, len(converted)), converted)

	return C.BitwardenItemSlice{items: cItems, len: C.size_t(len(converted))}
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
