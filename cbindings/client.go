package main

/*
#include "bw_common.h"

typedef struct {
	const char* baseUrl;
	const char* identityUrl;
	const char* apiUrl;
	const Handle* httpClient;
	const UUID* deviceId;
} NewClientOptions;
*/
import "C"
import (
	"net/http"

	"go.chrastecky.dev/bitsailor-core/bitwarden"
)

//export BitwardenNewClient
func BitwardenNewClient(outHandle *C.ClientHandle, options C.NewClientOptions) C.BitwardenResult {
	if outHandle == nil {
		setLastError(nullPointerError("outHandle"))
		return BitwardenError
	}

	goOptions := make([]bitwarden.Option, 0)
	if options.baseUrl != nil {
		goOptions = append(goOptions, bitwarden.WithBaseURL(C.GoString(options.baseUrl)))
	}
	if options.identityUrl != nil {
		goOptions = append(goOptions, bitwarden.WithIdentityURL(C.GoString(options.identityUrl)))
	}
	if options.apiUrl != nil {
		goOptions = append(goOptions, bitwarden.WithAPIURL(C.GoString(options.apiUrl)))
	}
	if options.httpClient != nil {
		httpClient, err := getHandleObj[*http.Client](handle(*options.httpClient))
		if err != nil {
			setLastError(err)
			return BitwardenError
		}

		goOptions = append(goOptions, bitwarden.WithHTTPClient(httpClient))
	}
	if options.deviceId != nil {
		goOptions = append(goOptions, bitwarden.WithDeviceID(parseUUIDFromC(*options.deviceId)))
	}

	client, err := bitwarden.NewClient(goOptions...)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	result := registerHandle(client)
	*outHandle = C.ClientHandle(result)

	clearLastError()
	return BitwardenSuccess
}
