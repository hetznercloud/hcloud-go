package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
		WithBackoffFunc(func(_ int) time.Duration { return 0 }),
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

func TestClientMeta(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("RateLimit-Limit", "1000")
		w.Header().Set("RateLimit-Remaining", "999")
		w.Header().Set("RateLimit-Reset", "1511954577")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"foo": "bar",
			"meta": {
				"pagination": {
					"page": 1
				}
			}
		}`)
	})

	ctx := context.Background()
	req, err := env.Client.NewRequest(ctx, "GET", "/", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	response, err := env.Client.Do(req, nil)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	if response.Meta.Ratelimit.Limit != 1000 {
		t.Errorf("unexpected ratelimit limit: %d", response.Meta.Ratelimit.Limit)
	}
	if response.Meta.Ratelimit.Remaining != 999 {
		t.Errorf("unexpected ratelimit remaining: %d", response.Meta.Ratelimit.Remaining)
	}
	if !response.Meta.Ratelimit.Reset.Equal(time.Unix(1511954577, 0)) {
		t.Errorf("unexpected ratelimit reset: %v", response.Meta.Ratelimit.Reset)
	}

	if response.Meta.Pagination.Page != 1 {
		t.Error("missing pagination")
	}
}

func TestClientMetaNonJSON(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "foo")
	})

	ctx := context.Background()
	req, err := env.Client.NewRequest(ctx, "GET", "/", nil)
	if err != nil {
		t.Fatalf("error creating request: %s", err)
	}

	response, err := env.Client.Do(req, nil)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	if response.Meta.Pagination != nil {
		t.Fatal("pagination should not be present")
	}
}

func TestResponseMetaPaginationUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		Data       []byte
		Pagination ResponseMetaPagination
	}{
		{
			Data: []byte(`{
				"page": 1,
				"per_page": 25,
				"previous_page": null,
				"next_page": null,
				"last_page": 1,
				"total_entries": 5
		        }`),
			Pagination: ResponseMetaPagination{
				Page:         1,
				PerPage:      25,
				PreviousPage: 0,
				NextPage:     0,
				LastPage:     1,
				TotalEntries: 5,
			},
		},
		{
			Data: []byte(`{
				"page": 2,
				"per_page": 25,
				"previous_page": 1,
				"next_page": 3,
				"last_page": 13,
				"total_entries": 322
		        }`),
			Pagination: ResponseMetaPagination{
				Page:         2,
				PerPage:      25,
				PreviousPage: 1,
				NextPage:     3,
				LastPage:     13,
				TotalEntries: 322,
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test case %d", i+1), func(t *testing.T) {
			var p ResponseMetaPagination
			if err := json.Unmarshal(testCase.Data, &p); err != nil {
				t.Fatal(err)
			}
			if p.Page != testCase.Pagination.Page {
				t.Errorf("unexpected page: %v", p.Page)
			}
			if p.PerPage != testCase.Pagination.PerPage {
				t.Errorf("unexpected per page: %v", p.PerPage)
			}
			if p.PreviousPage != testCase.Pagination.PreviousPage {
				t.Errorf("unexpected previous page: %v", p.PreviousPage)
			}
			if p.NextPage != testCase.Pagination.NextPage {
				t.Errorf("unexpected next page: %d", p.NextPage)
			}
			if p.LastPage != testCase.Pagination.LastPage {
				t.Errorf("unexpected last page: %d", p.LastPage)
			}
			if p.TotalEntries != testCase.Pagination.TotalEntries {
				t.Errorf("unexpected total entries: %d", p.TotalEntries)
			}
		})
	}
}
