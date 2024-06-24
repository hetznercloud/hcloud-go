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
		want    int
	}{
		{
			name: "network error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return nil, fmt.Errorf("network error")
			},
			want: 0,
		},
		{
			name: "http 503 error",
			wrapped: func(req *http.Request, _ any) (*Response, error) {
				resp := fakeResponse(t, 503, "", false)
				resp.Response.Request = req
				return resp, fmt.Errorf("%w %d", ErrStatusCode, 503)
			},
			want: 1,
		},
		{
			name: "api conflict error",
			wrapped: func(req *http.Request, _ any) (*Response, error) {
				resp := fakeResponse(t, 409, "", false)
				resp.Response.Request = req
				return nil, ErrorFromSchema(schema.Error{Code: string(ErrorCodeConflict)})
			},
			want: 1,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := &mockHandler{testCase.wrapped}

			retryCount := 0
			h := wrapRetryHandler(m, func(_ int) time.Duration {
				// Reset the mock handler to exit the retry loop
				m.f = func(_ *http.Request, _ any) (*Response, error) { return nil, nil }

				retryCount++
				return 0
			})

			client := NewClient(WithToken("dummy"))
			req, err := client.NewRequest(context.Background(), "GET", "/", nil)
			require.NoError(t, err)

			assert.Equal(t, 0, retryCount)

			h.Do(req, nil)

			assert.Equal(t, testCase.want, retryCount)
		})
	}
}
