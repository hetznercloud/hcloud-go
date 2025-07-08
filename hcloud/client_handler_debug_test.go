package hcloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebugHandler(t *testing.T) {
	testCases := []struct {
		name    string
		wrapped func(req *http.Request, v any) (*Response, error)
		want    string
	}{
		{
			name: "network error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return nil, fmt.Errorf("network error")
			},
			want: `--- Request:
GET /v1/ HTTP/1.1
Host: api.hetzner.cloud
User-Agent: hcloud-go/testing
Accept: application/json
Authorization: REDACTED
Accept-Encoding: gzip



`,
		},
		{
			name: "http 503 error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 503, "", false), nil
			},
			want: `--- Request:
GET /v1/ HTTP/1.1
Host: api.hetzner.cloud
User-Agent: hcloud-go/testing
Accept: application/json
Authorization: REDACTED
Accept-Encoding: gzip



--- Response:
HTTP/1.1 503 Service Unavailable
Connection: close



`,
		},
		{
			name: "http 200",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 200, `{"data": {"id": 1234, "name": "testing"}}`, true), nil
			},
			want: `--- Request:
GET /v1/ HTTP/1.1
Host: api.hetzner.cloud
User-Agent: hcloud-go/testing
Accept: application/json
Authorization: REDACTED
Accept-Encoding: gzip



--- Response:
HTTP/1.1 200 OK
Connection: close
Content-Type: application/json

{"data": {"id": 1234, "name": "testing"}}

`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			m := &mockHandler{testCase.wrapped}
			h := wrapDebugHandler(m, buf)

			client := NewClient(WithToken("dummy"))
			client.userAgent = "hcloud-go/testing"

			req, err := client.NewRequest(context.Background(), "GET", "/", nil)
			require.NoError(t, err)

			h.Do(req, nil)

			re := regexp.MustCompile(`\r`)
			output := re.ReplaceAllString(buf.String(), "")
			assert.Equal(t, testCase.want, output)
		})
	}
}
