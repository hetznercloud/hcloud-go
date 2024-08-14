package hcloud

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type debugWriteFunc func(id string, dump []byte)

func wrapDebugHandler(
	wrapped handler,
	opts *DebugOpts,
) handler {
	return &debugHandler{
		handler:       wrapped,
		writeRequest:  opts.WriteRequest,
		writeResponse: opts.WriteResponse,
	}
}

type debugHandler struct {
	handler       handler
	writeRequest  debugWriteFunc
	writeResponse debugWriteFunc
}

func (h *debugHandler) Do(req *http.Request, v any) (resp *Response, err error) {
	id := generateRandomID()

	// Clone the request, so we can redact the auth header, read the body
	// and use a new context.
	cloned, err := cloneRequest(req, context.Background())
	if err != nil {
		return nil, err
	}

	cloned.Header.Set("Authorization", "REDACTED")

	dumpReq, err := httputil.DumpRequestOut(cloned, true)
	if err != nil {
		return nil, err
	}

	h.writeRequest(id, dumpReq)

	resp, err = h.handler.Do(req, v)
	if err != nil {
		return resp, err
	}

	dumpResp, err := httputil.DumpResponse(resp.Response, true)
	if err != nil {
		return nil, err
	}

	h.writeResponse(id, dumpResp)

	return resp, err
}

func generateRandomID() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("failed to generate random string: %w", err))
	}
	return hex.EncodeToString(b)
}

func defaultDebugWriter(output io.Writer, title string) func(id string, dump []byte) {
	return func(_ string, dump []byte) {
		fmt.Fprintf(output,
			"--- %s:\n%s\n\n",
			title,
			bytes.Trim(dump, "\n"),
		)
	}
}
