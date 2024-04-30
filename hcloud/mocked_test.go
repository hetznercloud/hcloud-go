package hcloud

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockedTestCase struct {
	Name         string
	WantRequests []MockedRequest
	Run          func(env testEnv)
}

type MockedRequest struct {
	Method              string
	Path                string
	WantRequestBodyFunc func(t *testing.T, r *http.Request, body []byte)

	Status int
	Body   string
}

func RunMockedTestCases(t *testing.T, testCases []MockedTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			env := newTestEnvWithServer(httptest.NewServer(MockedRequestHandler(t, testCase.WantRequests)), nil)
			defer env.Teardown()

			testCase.Run(env)
		})
	}
}

func MockedRequestHandler(t *testing.T, requests []MockedRequest) http.HandlerFunc {
	index := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if testing.Verbose() {
			t.Logf("request %d: %s %s\n", index, r.Method, r.URL.Path)
		}

		if index >= len(requests) {
			t.Fatalf("received unknown request %d", index)
		}

		response := requests[index]
		assert.Equal(t, response.Method, r.Method)
		assert.Equal(t, response.Path, r.RequestURI)

		if response.WantRequestBodyFunc != nil {
			buffer, err := io.ReadAll(r.Body)
			defer func() {
				if err := r.Body.Close(); err != nil {
					t.Fatal(err)
				}
			}()
			if err != nil {
				t.Fatal(err)
			}
			response.WantRequestBodyFunc(t, r, buffer)
		}

		w.WriteHeader(response.Status)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(response.Body))
		if err != nil {
			t.Fatal(err)
		}

		index++
	})
}
