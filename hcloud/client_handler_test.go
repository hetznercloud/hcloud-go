package hcloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	f func(req *http.Request, v any) (resp *Response, err error)
}

func (h *mockHandler) Do(req *http.Request, v interface{}) (*Response, error) { return h.f(req, v) }

func fakeResponse(t *testing.T, statusCode int, body string, json bool) *Response {
	t.Helper()

	w := httptest.NewRecorder()
	if body != "" && json {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(statusCode)

	if body != "" {
		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}

	resp := &Response{Response: w.Result()} //nolint: bodyclose
	require.NoError(t, resp.populateBody())

	return resp
}

func TestCloneRequest(t *testing.T) {
	ctx := context.Background()

	req, err := http.NewRequest("GET", "/", bytes.NewBufferString("Hello"))
	require.NoError(t, err)
	req.Header.Set("Authorization", "sensitive")

	cloned, err := cloneRequest(req, ctx)
	require.NoError(t, err)
	cloned.Header.Set("Authorization", "REDACTED")
	cloned.Body = io.NopCloser(bytes.NewBufferString("Changed"))

	// Check context
	assert.Equal(t, req.Context(), cloned.Context())

	// Check headers
	assert.Equal(t, req.Header.Get("Authorization"), "sensitive")

	// Check body
	reqBody, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	assert.Equal(t, string(reqBody), "Hello")
}
