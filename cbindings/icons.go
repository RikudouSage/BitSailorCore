package main

/*
#include "bw_common.h"
#include "bw_icons.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

//export BitwardenGetIcon
func BitwardenGetIcon(
	client C.ClientHandle,
	ctx C.ContextHandle,
	hostname *C.char,
	out *C.BitwardenByteSlice,
) C.BitwardenResult {
	if hostname == nil {
		setLastError(nullPointerError("hostname"))
		return BitwardenError
	}
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	iconsGo, ctxGo, err := getCommonIconsHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	icon, err := iconsGo.GetIcon(ctxGo, C.GoString(hostname))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	bitwardenByteSliceIntoC(out, icon)
	clearLastError()
	return BitwardenSuccess
}

//export BitwardenGetIcons
func BitwardenGetIcons(
	client C.ClientHandle,
	ctx C.ContextHandle,
	hostnames C.BitwardenStringSlice,
	out *C.BitwardenIconSlice,
) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	iconsGo, ctxGo, err := getCommonIconsHandles(client, ctx)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	icons, err := iconsGo.GetIcons(ctxGo, goStringSliceFromC(hostnames))
	bitwardenIconSliceIntoC(out, icons)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenFreeByteSlice
func BitwardenFreeByteSlice(value *C.BitwardenByteSlice) {
	freeBitwardenByteSlice(value)
}

//export BitwardenFreeIcons
func BitwardenFreeIcons(value *C.BitwardenIconSlice) {
	freeBitwardenIconSlice(value)
}

func bitwardenByteSliceIntoC(out *C.BitwardenByteSlice, value []byte) {
	if len(value) == 0 {
		clearC(out)
		return
	}

	items := C.CBytes(value)
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.items), items)
	putCValue(unsafe.Pointer(out), unsafe.Offsetof(out.len), C.size_t(len(value)))
}

func freeBitwardenByteSlice(value *C.BitwardenByteSlice) {
	if value == nil {
		return
	}

	C.free(unsafe.Pointer(value.items))
	clearC(value)
}

func bitwardenIconSliceIntoC(out *C.BitwardenIconSlice, icons map[string][]byte) {
	if len(icons) == 0 {
		clearC(out)
		return
	}

	items := (*C.BitwardenIcon)(C.malloc(C.size_t(len(icons)) * C.size_t(unsafe.Sizeof(C.BitwardenIcon{}))))
	cIcons := unsafe.Slice(items, len(icons))
	index := 0
	for hostname, icon := range icons {
		cIcons[index].hostname = C.CString(hostname)
		bitwardenByteSliceIntoC(&cIcons[index].icon, icon)
		index++
	}

	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.items), unsafe.Pointer(items))
	putCValue(unsafe.Pointer(out), unsafe.Offsetof(out.len), C.size_t(len(icons)))
}

func freeBitwardenIconSlice(value *C.BitwardenIconSlice) {
	if value == nil {
		return
	}

	icons := unsafe.Slice(value.items, int(value.len))
	for i := range icons {
		C.free(unsafe.Pointer(icons[i].hostname))
		freeBitwardenByteSlice(&icons[i].icon)
	}
	C.free(unsafe.Pointer(value.items))

	clearC(value)
}
