package mockutils

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Request describes a http request that a [httptest.Server] should receive, and the
// corresponding response to return.
//
// Additional checks on the request (e.g. request body) may be added using the
// [Request.Want] function.
//
// The response body is populated from either a JSON struct, or a JSON string.
type Request struct {
	Method string
	Path   string
	Want   func(t *testing.T, r *http.Request)

	Status  int
	JSON    any
	JSONRaw string
}

// Handler is used with a [httptest.Server] to mock http requests provided by the user.
//
// Request matching is based on the request count, and the user provided request will be
// iterated over.
func Handler(t *testing.T, requests []Request) http.HandlerFunc {
	t.Helper()

	index := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if testing.Verbose() {
			t.Logf("call %d: %s %s\n", index, r.Method, r.RequestURI)
		}

		if index >= len(requests) {
			t.Fatalf("received unknown call %d", index)
		}

		expected := requests[index]

		expectedCall := expected.Method
		foundCall := r.Method
		if expected.Path != "" {
			expectedCall += " " + expected.Path
			foundCall += " " + r.RequestURI
		}
		require.Equal(t, expectedCall, foundCall)

		if expected.Want != nil {
			expected.Want(t, r)
		}

		switch {
		case expected.JSON != nil:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(expected.Status)
			if err := json.NewEncoder(w).Encode(expected.JSON); err != nil {
				t.Fatal(err)
			}
		case expected.JSONRaw != "":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(expected.Status)
			_, err := w.Write([]byte(expected.JSONRaw))
			if err != nil {
				t.Fatal(err)
			}
		default:
			w.WriteHeader(expected.Status)
		}

		index++
	})
}
