package main

/*
#include "bw_common.h"
#include "bw_generator.h"
*/
import "C"
import "go.chrastecky.dev/bitwarden-client/bitwarden"

//export BitwardenGeneratePassword
func BitwardenGeneratePassword(client C.ClientHandle, request C.BitwardenPasswordGeneratorRequest, out **C.char) C.BitwardenResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return BitwardenError
	}
	clientGo, err := getHandleObj[bitwarden.Client](handle(client))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	requestGo := &bitwarden.PasswordGeneratorRequest{
		Lowercase:      goBoolFromCPtr(request.lowercase),
		Uppercase:      goBoolFromCPtr(request.uppercase),
		Numbers:        goBoolFromCPtr(request.numbers),
		Special:        goBoolFromCPtr(request.special),
		Length:         goIntFromCPtr(request.length),
		AvoidAmbiguous: goBoolFromCPtr(request.avoidAmbiguous),
		MinLowercase:   goIntFromCPtr(request.minLowercase),
		MinUppercase:   goIntFromCPtr(request.minUppercase),
		MinNumber:      goIntFromCPtr(request.minNumber),
		MinSpecial:     goIntFromCPtr(request.minSpecial),
	}

	password, err := clientGo.GeneratePassword(requestGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*out = C.CString(password)
	clearLastError()
	return BitwardenSuccess
}
