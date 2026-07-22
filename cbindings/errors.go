package main

/*
#include <bw_errors.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

func clearLastError() {
	C.bitwarden_clear_last_error()
}

func setLastError(err error) {
	if err == nil {
		clearLastError()
		return
	}
	setLastErrorMessage(err.Error())
}

func setLastErrorMessage(msg string) {
	if msg == "" {
		clearLastError()
		return
	}
	cstr := C.CString(msg)
	C.bitwarden_set_last_error_copy(cstr)
	C.free(unsafe.Pointer(cstr))
}

func nullPointerError(name string) error {
	return errors.New(name + " is NULL")
}

//export BitwardenGetLastError
func BitwardenGetLastError(buf *C.char, bufLen C.size_t) C.size_t {
	return C.bitwarden_get_last_error(buf, bufLen)
}
