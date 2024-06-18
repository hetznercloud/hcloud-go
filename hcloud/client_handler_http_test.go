package hcloud

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPHandler(t *testing.T) {
	testServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		}),
	)

	h := newHTTPHandler(&http.Client{})

	req, err := http.NewRequest("GET", testServer.URL, nil)
	require.NoError(t, err)

	resp, err := h.Do(req, nil)
	require.NoError(t, err)

	// Ensure the internal response body is populated
	assert.Equal(t, []byte("hello"), resp.body)

	// Ensure the original response body is readable by external users
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), body)
}
