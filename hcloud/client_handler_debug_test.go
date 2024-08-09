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
	generateTimestampOrig := generateTimestamp
	generateTimestamp = func() string { return "2024-08-08T12:15:18+02:00" }
	defer func() { generateTimestamp = generateTimestampOrig }()

	generateRandomIDOrig := generateRandomID
	generateRandomID = func() string { return "22ae0311" }
	defer func() { generateRandomID = generateRandomIDOrig }()

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
			want: `2024-08-08T12:15:18+02:00 [22ae0311]: --- Request:
2024-08-08T12:15:18+02:00 [22ae0311]: GET /v1/ HTTP/1.1
2024-08-08T12:15:18+02:00 [22ae0311]: Host: api.hetzner.cloud
2024-08-08T12:15:18+02:00 [22ae0311]: User-Agent: hcloud-go/testing
2024-08-08T12:15:18+02:00 [22ae0311]: Authorization: REDACTED
2024-08-08T12:15:18+02:00 [22ae0311]: Accept-Encoding: gzip
2024-08-08T12:15:18+02:00 [22ae0311]: 
2024-08-08T12:15:18+02:00 [22ae0311]: 
`,
		},
		{
			name: "http 503 error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 503, "", false), nil
			},
			want: `2024-08-08T12:15:18+02:00 [22ae0311]: --- Request:
2024-08-08T12:15:18+02:00 [22ae0311]: GET /v1/ HTTP/1.1
2024-08-08T12:15:18+02:00 [22ae0311]: Host: api.hetzner.cloud
2024-08-08T12:15:18+02:00 [22ae0311]: User-Agent: hcloud-go/testing
2024-08-08T12:15:18+02:00 [22ae0311]: Authorization: REDACTED
2024-08-08T12:15:18+02:00 [22ae0311]: Accept-Encoding: gzip
2024-08-08T12:15:18+02:00 [22ae0311]: 
2024-08-08T12:15:18+02:00 [22ae0311]: 
2024-08-08T12:15:18+02:00 [22ae0311]: --- Response:
2024-08-08T12:15:18+02:00 [22ae0311]: HTTP/1.1 503 Service Unavailable
2024-08-08T12:15:18+02:00 [22ae0311]: Connection: close
2024-08-08T12:15:18+02:00 [22ae0311]: 
2024-08-08T12:15:18+02:00 [22ae0311]: 
`,
		},
		{
			name: "http 200",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 200, `{"data": {"id": 1234, "name": "testing"}}`, true), nil
			},
			want: `2024-08-08T12:15:18+02:00 [22ae0311]: --- Request:
2024-08-08T12:15:18+02:00 [22ae0311]: GET /v1/ HTTP/1.1
2024-08-08T12:15:18+02:00 [22ae0311]: Host: api.hetzner.cloud
2024-08-08T12:15:18+02:00 [22ae0311]: User-Agent: hcloud-go/testing
2024-08-08T12:15:18+02:00 [22ae0311]: Authorization: REDACTED
2024-08-08T12:15:18+02:00 [22ae0311]: Accept-Encoding: gzip
2024-08-08T12:15:18+02:00 [22ae0311]: 
2024-08-08T12:15:18+02:00 [22ae0311]: 
2024-08-08T12:15:18+02:00 [22ae0311]: --- Response:
2024-08-08T12:15:18+02:00 [22ae0311]: HTTP/1.1 200 OK
2024-08-08T12:15:18+02:00 [22ae0311]: Connection: close
2024-08-08T12:15:18+02:00 [22ae0311]: Content-Type: application/json
2024-08-08T12:15:18+02:00 [22ae0311]: 
2024-08-08T12:15:18+02:00 [22ae0311]: {"data": {"id": 1234, "name": "testing"}}
2024-08-08T12:15:18+02:00 [22ae0311]: 
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
