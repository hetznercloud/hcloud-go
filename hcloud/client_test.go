package hcloud

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testEnv struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	Client *Client
}

func (env *testEnv) Teardown() {
	env.Server.Close()
	env.Server = nil
	env.Mux = nil
	env.Client = nil
}

func newTestEnv() testEnv {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(
		WithEndpoint(server.URL),
		WithToken("token"),
	)
	return testEnv{
		Server: server,
		Mux:    mux,
		Client: client,
	}
}

func TestClientEndpointTrailingSlashesRemoved(t *testing.T) {
	client := NewClient(WithEndpoint("http://api/v1.0/////"))
	if strings.HasSuffix(client.endpoint, "/") {
		t.Fatalf("endpoint has trailing slashes: %q", client.endpoint)
	}
}

func TestClientError(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{
			"error": {
				"code": "service_error",
				"message": "An error occured"
			}
		}`)
	})

	ctx := context.Background()
	req, err := env.Client.NewRequest(ctx, "GET", "/error", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	_, err = env.Client.Do(req, nil)
	if _, ok := err.(Error); !ok {
		t.Fatalf("unexpected error of type %T: %v", err, err)
	}

	apiError := err.(Error)

	if apiError.Code != "service_error" {
		t.Errorf("unexpected error code: %q", apiError.Code)
	}
	if apiError.Message != "An error occured" {
		t.Errorf("unexpected error message: %q", apiError.Message)
	}
}
