package hcloud

import (
	"net/http"
	"strconv"
	"time"
)

func wrapRateLimitHandler(wrapped handler) handler {
	return &rateLimitHandler{wrapped}
}

type rateLimitHandler struct {
	handler handler
}

func (h *rateLimitHandler) Do(req *http.Request, v any) (resp *Response, err error) {
	resp, err = h.handler.Do(req, v)

	if resp != nil && resp.Header != nil {
		if h := resp.Header.Get("RateLimit-Limit"); h != "" {
			resp.Meta.Ratelimit.Limit, _ = strconv.Atoi(h)
		}
		if h := resp.Header.Get("RateLimit-Remaining"); h != "" {
			resp.Meta.Ratelimit.Remaining, _ = strconv.Atoi(h)
		}
		if h := resp.Header.Get("RateLimit-Reset"); h != "" {
			if ts, err := strconv.ParseInt(h, 10, 64); err == nil {
				resp.Meta.Ratelimit.Reset = time.Unix(ts, 0)
			}
		}
	}

	return resp, err
}
