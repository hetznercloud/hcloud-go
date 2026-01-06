package hcloud

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestRetryHandler(t *testing.T) {
	testCases := []struct {
		name    string
		wrapped func(req *http.Request, v any) (*Response, error)
		recover bool
		want    func(t *testing.T, err error, retryCount int)
	}{
		{
			name: "random error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return nil, fmt.Errorf("random error")
			},
			want: func(t *testing.T, err error, retryCount int) {
				assert.EqualError(t, err, "random error")
				assert.Equal(t, 0, retryCount)
			},
		},
		{
			name: "http 502 error recovery",
			wrapped: func(req *http.Request, _ any) (*Response, error) {
				resp := fakeResponse(t, 502, "", false)
				resp.Response.Request = req
				return resp, fmt.Errorf("%w %d", ErrStatusCode, 502)
			},
			recover: true,
			want: func(t *testing.T, err error, retryCount int) {
				assert.NoError(t, err)
				assert.Equal(t, 1, retryCount)
			},
		},
		{
			name: "http 502 error",
			wrapped: func(req *http.Request, _ any) (*Response, error) {
				resp := fakeResponse(t, 502, "", false)
				resp.Response.Request = req
				return resp, fmt.Errorf("%w %d", ErrStatusCode, 502)
			},
			want: func(t *testing.T, err error, retryCount int) {
				assert.EqualError(t, err, "server responded with status code 502")
				assert.Equal(t, 5, retryCount)
			},
		},
		{
			name: "api conflict error recovery",
			wrapped: func(req *http.Request, _ any) (*Response, error) {
				resp := fakeResponse(t, 409, "", false)
				resp.Response.Request = req
				return nil, ErrorFromSchema(schema.Error{Code: string(ErrorCodeConflict), Message: "A conflict occurred"})
			},
			recover: true,
			want: func(t *testing.T, err error, retryCount int) {
				assert.NoError(t, err)
				assert.Equal(t, 1, retryCount)
			},
		},
		{
			name: "api conflict error",
			wrapped: func(req *http.Request, _ any) (*Response, error) {
				resp := fakeResponse(t, 409, "", false)
				resp.Response.Request = req
				return nil, ErrorFromSchema(schema.Error{Code: string(ErrorCodeConflict), Message: "A conflict occurred"})
			},
			want: func(t *testing.T, err error, retryCount int) {
				assert.EqualError(t, err, "A conflict occurred (conflict)")
				assert.Equal(t, 5, retryCount)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := &mockHandler{testCase.wrapped}

			retryCount := 0
			h := wrapRetryHandler(m, func(retries int) time.Duration {
				assert.Equal(t, retryCount, retries)

				if testCase.recover {
					// Reset the mock handler to exit the retry loop
					m.f = func(_ *http.Request, _ any) (*Response, error) { return nil, nil }
				}

				retryCount++
				return 0
			}, 5)

			client := NewClient(WithToken("dummy"))
			req, err := client.NewRequest(context.Background(), "GET", "/", nil)
			require.NoError(t, err)
			require.Equal(t, 0, retryCount)

			_, err = h.Do(req, nil)
			testCase.want(t, err, retryCount)
		})
	}
}

func TestRetryPolicy(t *testing.T) {
	testCases := []struct {
		name string
		resp *Response
		want bool
	}{
		{
			name: "server returns 502 error",
			resp: fakeResponse(t, 502, ``, false),
			want: true,
		},
		{
			name: "api returns unavailable error",
			resp: fakeResponse(t, 503, `{"error":{"code":"unavailable"}}`, true),
			want: false,
		},
		{
			name: "server returns 503 error",
			resp: fakeResponse(t, 503, ``, false),
			want: false,
		},
		{
			name: "server returns timeout error",
			resp: fakeResponse(t, 504, `{"error":{"code":"timeout"}}`, true),
			want: true,
		},
		{
			name: "api returns rate_limit_exceeded error",
			resp: fakeResponse(t, 429, `{"error":{"code":"rate_limit_exceeded"}}`, true),
			want: true,
		},
		{
			name: "server returns 429 error",
			resp: fakeResponse(t, 429, ``, false),
			want: false,
		},
		{
			name: "api returns conflict error",
			resp: fakeResponse(t, 409, `{"error":{"code":"conflict"}}`, true),
			want: true,
		},
		{
			// HTTP 409 is used in many situations (e.g. uniqueness_error), we must only
			// retry if the API error code is conflict.
			name: "server returns 409 error",
			resp: fakeResponse(t, 409, ``, false),
			want: false,
		},
		{
			// The API error code locked is used in many unexpected situations, we can
			// only retry in specific context where we know the error is not misused.
			name: "api returns locked error",
			resp: fakeResponse(t, 423, `{"error":{"code":"locked"}}`, true),
			want: false,
		},
		{
			// HTTP 423 is used in many situations (e.g. protected), we must only
			// retry if the API error code is locked.
			name: "server returns 423 error",
			resp: fakeResponse(t, 423, ``, false),
			want: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)

			m := &mockHandler{func(req *http.Request, _ any) (*Response, error) {
				testCase.resp.Request = req
				return testCase.resp, nil
			}}
			h := wrapErrorHandler(m)

			result := retryPolicy(h.Do(req, nil))
			assert.Equal(t, testCase.want, result)
		})
	}
}
