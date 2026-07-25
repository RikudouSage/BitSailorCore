package main

/*
#include "bw_common.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"

	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenExportSession
func BitwardenExportSession(session C.SessionHandle, out **C.char) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}

	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sessionBytes, err := json.Marshal(sessionGo)
	if err != nil {
		setLastError(fmt.Errorf("failed encoding session: %w", err))
		return BitwardenError
	}

	*out = C.CString(string(sessionBytes))
	clearLastError()
	return BitwardenSuccess
}
