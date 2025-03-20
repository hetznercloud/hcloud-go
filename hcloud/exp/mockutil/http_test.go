package mockutil

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	server := NewServer(t, []Request{
		{
			Method: "GET", Path: "/",
			Status: 200,
			JSON: struct {
				Data string `json:"data"`
			}{
				Data: "Hello",
			},
		},
		{
			Method: "GET", Path: "/",
			Status:  400,
			JSONRaw: `{"error": "failed"}`,
		},
		{
			Method: "GET", Path: "/",
			Status: 503,
		},
		{
			Method: "GET",
			Want: func(t *testing.T, r *http.Request) {
				require.True(t, strings.HasPrefix(r.RequestURI, "/random?key="))
			},
			Status: 200,
		},
	})

	// Request 1
	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.JSONEq(t, `{"data":"Hello"}`, readBody(t, resp))

	// Request 2
	resp, err = http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.JSONEq(t, `{"error": "failed"}`, readBody(t, resp))

	// Request 3
	resp, err = http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)
	assert.Equal(t, "", resp.Header.Get("Content-Type"))
	assert.Equal(t, "", readBody(t, resp))

	// Request 4
	resp, err = http.Get(fmt.Sprintf("%s/random?key=%d", server.URL, rand.Int63()))
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "", resp.Header.Get("Content-Type"))
	assert.Equal(t, "", readBody(t, resp))

	// Extra request 5
	server.Expect([]Request{
		{Method: "GET", Path: "/", Status: 200},
	})

	resp, err = http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "", resp.Header.Get("Content-Type"))
	assert.Equal(t, "", readBody(t, resp))
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return strings.TrimSuffix(string(body), "\n")
}
