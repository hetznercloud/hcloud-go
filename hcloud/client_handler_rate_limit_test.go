package hcloud

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitHandler(t *testing.T) {
	testCases := []struct {
		name    string
		wrapped func(req *http.Request, v any) (*Response, error)
		want    func(t *testing.T, resp *Response, err error)
	}{
		{
			name: "response",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				resp := fakeResponse(t, 200, "", false)
				resp.Header.Set("RateLimit-Limit", "1000")
				resp.Header.Set("RateLimit-Remaining", "999")
				resp.Header.Set("RateLimit-Reset", "1511954577")
				return resp, nil
			},
			want: func(t *testing.T, resp *Response, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1000, resp.Meta.Ratelimit.Limit)
				assert.Equal(t, 999, resp.Meta.Ratelimit.Remaining)
				assert.Equal(t, time.Unix(1511954577, 0), resp.Meta.Ratelimit.Reset)
			},
		},
		{
			name: "any error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return nil, fmt.Errorf("any error")
			},
			want: func(t *testing.T, resp *Response, err error) {
				assert.EqualError(t, err, "any error")
				assert.Nil(t, resp)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := &mockHandler{testCase.wrapped}
			h := wrapRateLimitHandler(m)

			resp, err := h.Do(nil, nil)

			testCase.want(t, resp, err)
		})
	}
}
