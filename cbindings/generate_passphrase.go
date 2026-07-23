package main

/*
#include "bw_common.h"
#include "bw_generator.h"
*/
import "C"
import "go.chrastecky.dev/bitwarden-client/bitwarden"

//export BitwardenGeneratePassphrase
func BitwardenGeneratePassphrase(client C.ClientHandle, request C.BitwardenPassphraseGeneratorRequest, out **C.char) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}
	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	requestGo := &bitwarden.PassphraseGeneratorRequest{
		NumWords:      goIntFromCPtr(request.numWords),
		WordSeparator: goStringFromCPtr(request.wordSeparator),
		Capitalize:    goBoolFromCPtr(request.capitalize),
		IncludeNumber: goBoolFromCPtr(request.includeNumber),
	}

	passphrase, err := clientGo.GeneratePassphrase(requestGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*out = C.CString(passphrase)
	clearLastError()
	return BitwardenSuccess
}
