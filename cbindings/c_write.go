package main

/*
#include <stdlib.h>
#include <stdbool.h>
*/
import "C"
import "unsafe"

func clearC[T any](value *T) {
	if value == nil {
		return
	}

	clear(unsafe.Slice((*byte)(unsafe.Pointer(value)), unsafe.Sizeof(*value)))
}

func putCValue[T any](base unsafe.Pointer, offset uintptr, value T) {
	*(*T)(unsafe.Add(base, offset)) = value
}

func putCPtr(base unsafe.Pointer, offset uintptr, value unsafe.Pointer) {
	*(*uintptr)(unsafe.Add(base, offset)) = uintptr(value)
}

func putCSlice(base unsafe.Pointer, sliceOffset, itemsOffset, lenOffset uintptr, items unsafe.Pointer, length C.size_t) {
	sliceBase := unsafe.Add(base, sliceOffset)
	putCPtr(sliceBase, itemsOffset, items)
	putCValue(sliceBase, lenOffset, length)
}

func cPtr[T any](value *T) unsafe.Pointer {
	return unsafe.Pointer(value)
}

func cStringPtr(value string) unsafe.Pointer {
	return unsafe.Pointer(C.CString(value))
}

func cStringPtrFromPtr(value *string) unsafe.Pointer {
	if value == nil {
		return nil
	}
	return cStringPtr(*value)
}

func cBoolPtr(value *bool) unsafe.Pointer {
	if value == nil {
		return nil
	}

	out := (*C.bool)(C.malloc(C.size_t(unsafe.Sizeof(C.bool(false)))))
	*out = C.bool(*value)
	return unsafe.Pointer(out)
}

func cIntPtr(value *int) unsafe.Pointer {
	if value == nil {
		return nil
	}

	out := (*C.int)(C.malloc(C.size_t(unsafe.Sizeof(C.int(0)))))
	*out = C.int(*value)
	return unsafe.Pointer(out)
}
