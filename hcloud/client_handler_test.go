package hcloud

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
