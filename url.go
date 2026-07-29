package bitwarden

import (
	"net/url"
	"strings"
)

func urlWithPath(base *url.URL, endpointPath string) *url.URL {
	uri := new(*base)
	basePath := strings.TrimRight(uri.Path, "/")
	endpointPath = "/" + strings.TrimLeft(endpointPath, "/")
	if basePath == "" {
		uri.Path = endpointPath
	} else {
		uri.Path = basePath + endpointPath
	}
	uri.RawPath = ""

	return uri
}
