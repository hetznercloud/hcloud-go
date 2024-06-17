package hcloud

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func wrapDebugHandler(wrapped handler, output io.Writer) handler {
	return &debugHandler{wrapped, output}
}

type debugHandler struct {
	handler handler
	output  io.Writer
}

func (h *debugHandler) Do(req *http.Request, v any) (resp *Response, err error) {
	// Duplicate the request, so we can redact the auth header and read the body

	// Clone the request using the new context
	cloned := req.Clone(context.Background())

	if req.ContentLength > 0 {
		cloned.Body, err = req.GetBody()
		if err != nil {
			return nil, err
		}
	}

	cloned.Header.Set("Authorization", "REDACTED")

	dumpReq, err := httputil.DumpRequestOut(cloned, true)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(h.output, "--- Request:\n%s\n\n", dumpReq)

	resp, err = h.handler.Do(req, v)
	if err != nil {
		return resp, err
	}

	dumpResp, err := httputil.DumpResponse(resp.Response, true)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintf(h.output, "--- Response:\n%s\n\n", dumpResp)

	return resp, err
}
