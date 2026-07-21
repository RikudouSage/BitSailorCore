package http

import (
	"io"
	"net/http"
)

func DrainResponse(response *http.Response) {
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()
}
