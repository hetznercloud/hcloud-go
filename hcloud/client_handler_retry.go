package hcloud

import (
	"errors"
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
	ctx := req.Context()

	for {
		// Clone the request using the original context
		cloned, err := cloneRequest(req, ctx)
		if err != nil {
			return nil, err
		}

		resp, err = h.handler.Do(cloned, v)
		if err != nil {
			// Beware the diversity of the errors:
			// - request preparation
			// - network connectivity
			// - http status code (see [errorHandler])
			if ctx.Err() != nil {
				return resp, ctx.Err()
			}

			if retryPolicy(resp, err) {
				select {
				case <-ctx.Done():
					return resp, err
				case <-time.After(h.backoffFunc(retries)):
					retries++
					continue
				}
			}
		}

		return resp, err
	}
}

func retryPolicy(resp *Response, err error) bool {
	if err != nil {
		var apiErr Error

		switch {
		case errors.As(err, &apiErr):
			if apiErr.Code == ErrorCodeConflict {
				return true
			}
		case errors.Is(err, ErrorStatusCode):
			if resp == nil || resp.Response == nil {
				// Should not happen
				return false
			}

			if resp.Response.Request.Method != "GET" {
				return false
			}

			switch resp.Response.StatusCode {
			// 4xx errors
			case http.StatusTooManyRequests:
				return true
			// 5xx errors
			case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
				return true
			}
		}
	}

	return false
}
