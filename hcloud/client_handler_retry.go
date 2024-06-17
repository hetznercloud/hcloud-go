package hcloud

import (
	"net/http"
	"time"
)

func wrapRetryHandler(wrapped handler, backoffFunc BackoffFunc) handler {
	return &retryHandler{wrapped, backoffFunc}
}

type retryHandler struct {
	handler     handler
	backoffFunc BackoffFunc
}

func (h *retryHandler) Do(req *http.Request, v any) (resp *Response, err error) {
	retries := 0

	for {
		// Clone the request using the original context
		cloned := req.Clone(req.Context())

		if req.ContentLength > 0 {
			cloned.Body, err = req.GetBody()
			if err != nil {
				return nil, err
			}
		}

		resp, err = h.handler.Do(cloned, v)
		if err != nil {
			// Beware the diversity of the errors:
			// - request preparation
			// - network connectivity
			// - http status code (see [errorHandler])
			// - response parsing
			if IsError(err, ErrorCodeConflict) {
				time.Sleep(h.backoffFunc(retries))
				retries++
				continue
			}
		}

		return resp, err
	}
}
