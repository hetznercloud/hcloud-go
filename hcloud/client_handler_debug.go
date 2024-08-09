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
	"strings"
	"time"
)

func wrapDebugHandler(wrapped handler, output io.Writer) handler {
	return &debugHandler{wrapped, output}
}

type debugHandler struct {
	handler handler
	output  io.Writer
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

	h.write(id, "Request", dumpReq)

	resp, err = h.handler.Do(req, v)
	if err != nil {
		return resp, err
	}

	dumpResp, err := httputil.DumpResponse(resp.Response, true)
	if err != nil {
		return nil, err
	}

	h.write(id, "Response", dumpResp)

	return resp, err
}

var generateRandomID = func() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("failed to generate random string: %w", err))
	}
	return hex.EncodeToString(b)
}

var generateTimestamp = func() string {
	return time.Now().Format(time.RFC3339)
}

func prependPrefix(prefix, input string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("%s: %s", prefix, line)
	}
	return strings.Join(lines, "\n")
}

func (h *debugHandler) write(id, title string, content []byte) {
	fmt.Fprintln(h.output,
		prependPrefix(
			fmt.Sprintf("%s [%s]", generateTimestamp(), id),
			fmt.Sprintf("--- %s:\n%s\n", title, bytes.Trim(content, "\n")),
		),
	)
}
