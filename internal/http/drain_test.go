package http

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type trackingReadCloser struct {
	*strings.Reader
	closed bool
}

func (receiver *trackingReadCloser) Close() error {
	receiver.closed = true
	return nil
}

func TestDrainResponseReadsAndClosesBody(t *testing.T) {
	t.Parallel()

	body := &trackingReadCloser{Reader: strings.NewReader("response body")}
	response := &http.Response{Body: body}

	DrainResponse(response)

	if !body.closed {
		t.Fatal("body.closed = false, want true")
	}
	remaining, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("ReadAll() returned error: %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("remaining body = %q, want empty", remaining)
	}
}
