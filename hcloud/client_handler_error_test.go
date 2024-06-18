package hcloud

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	testCases := []struct {
		name    string
		wrapped func(req *http.Request, v any) (*Response, error)
		want    func(t *testing.T, resp *Response, err error)
	}{
		{
			name: "no error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 200, `{"data": "Hello"}`, true), nil
			},
			want: func(t *testing.T, resp *Response, err error) {
				assert.Equal(t, 200, resp.StatusCode)
				assert.NoError(t, err)
			},
		},
		{
			name: "network error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return nil, fmt.Errorf("network error")
			},
			want: func(t *testing.T, resp *Response, err error) {
				assert.Nil(t, resp)
				assert.EqualError(t, err, "network error")
			},
		},
		{
			name: "http 503 error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 503, "", false), nil
			},
			want: func(t *testing.T, resp *Response, err error) {
				assert.Equal(t, 503, resp.StatusCode)
				assert.EqualError(t, err, "hcloud: server responded with status code 503")
			},
		},
		{
			name: "http 422 error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 422, `{"error": {"code": "service_error", "message": "An error occurred"}}`, true), nil
			},
			want: func(t *testing.T, resp *Response, err error) {
				assert.Equal(t, 422, resp.StatusCode)
				assert.EqualError(t, err, "An error occurred (service_error)")
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := &mockHandler{testCase.wrapped}
			h := wrapErrorHandler(m)

			resp, err := h.Do(nil, nil)

			testCase.want(t, resp, err)
		})
	}
}
