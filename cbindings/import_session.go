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

//export BitwardenImportSession
func BitwardenImportSession(inSession *C.SessionHandle, exportData *C.char, outSession *C.SessionHandle) C.BitwardenResult {
	if outSession == nil {
		setLastError(nullPointerError("outSession"))
		return BitwardenError
	}
	if exportData == nil {
		setLastError(nullPointerError("exportData"))
		return BitwardenError
	}

	var sessionGo *result.Session
	if inSession != nil {
		var err error
		sessionGo, err = getHandleObj[*result.Session](handle(*inSession))
		if err != nil {
			setLastError(err)
			return BitwardenError
		}
	}

	if err := json.Unmarshal([]byte(C.GoString(exportData)), &sessionGo); err != nil {
		setLastError(fmt.Errorf("failed decoding json data: %w", err))
		return BitwardenError
	}

	if inSession != nil {
		*outSession = *inSession
	} else {
		*outSession = C.SessionHandle(registerHandle(sessionGo))
	}

	return BitwardenSuccess
}
