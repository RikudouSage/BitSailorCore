package main

/*
#include <bw_common.h>
*/
import "C"

//export BitwardenCloseHandle
func BitwardenCloseHandle(handleID C.Handle) C.Result {
	err := unregisterHandle(handle(handleID))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	clearLastError()
	return BitwardenSuccess
}
