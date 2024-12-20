package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func makeTestUtils(t *testing.T) (context.Context, *mockutil.Server, *Client) {
	ctx := context.Background()

	server := mockutil.NewServer(t, nil)

	client := NewClient(
		WithEndpoint(server.URL),
		WithRetryOpts(RetryOpts{BackoffFunc: ConstantBackoff(0), MaxRetries: 5}),
		WithPollOpts(PollOpts{BackoffFunc: ConstantBackoff(0)}),
	)

	return ctx, server, client
}

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
	return newTestEnvWithServer(server, mux)
}

func newTestEnvWithServer(server *httptest.Server, mux *http.ServeMux) testEnv {
	client := NewClient(
		WithEndpoint(server.URL),
		WithToken("token"),
		WithRetryOpts(RetryOpts{
			BackoffFunc: ConstantBackoff(0),
			MaxRetries:  5,
		}),
		WithPollOpts(PollOpts{
			BackoffFunc: ConstantBackoff(0),
		}),
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
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code:    "service_error",
				Message: "An error occurred",
			},
		})
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
	if apiError.Message != "An error occurred" {
		t.Errorf("unexpected error message: %q", apiError.Message)
	}
	if apiError.Response().StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("unexpected http status code: %q", apiError.Response().StatusCode)
	}
}

func TestClientInvalidToken(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Client = NewClient(
		WithEndpoint(env.Server.URL),
		WithToken("invalid token\n"),
	)

	ctx := context.Background()
	_, err := env.Client.NewRequest(ctx, "GET", "/", nil)

	if nil == err {
		t.Error("Failed to trigger expected error")
	} else if err.Error() != "Authorization token contains invalid characters" {
		t.Fatalf("Invalid encoded authorization token triggered unexpected error message: %s", err)
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

func TestClientAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	var (
		ctx          = context.Background()
		conflicting  bool
		expectedPage = 1
	)

	env.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		respBody := schema.MetaResponse{
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					LastPage:     3,
					PerPage:      1,
					TotalEntries: 3,
				},
			},
		}

		switch page := r.URL.Query().Get("page"); page {
		case "", "1":
			respBody.Meta.Pagination.Page = 1
			respBody.Meta.Pagination.NextPage = 2
		case "2":
			if !conflicting {
				conflicting = true
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(schema.ErrorResponse{
					Error: schema.Error{
						Code:    string(ErrorCodeConflict),
						Message: "conflict",
					},
				})
				return
			}
			respBody.Meta.Pagination.Page = 2
			respBody.Meta.Pagination.PreviousPage = 1
			respBody.Meta.Pagination.NextPage = 3
		case "3":
			respBody.Meta.Pagination.Page = 3
			respBody.Meta.Pagination.PreviousPage = 2
		default:
			t.Errorf("bad page: %q", page)
		}

		json.NewEncoder(w).Encode(respBody)
	})

	_ = env.Client.all(func(page int) (*Response, error) {
		if page != expectedPage {
			t.Fatalf("expected page %d, but called for %d", expectedPage, page)
		}

		path := fmt.Sprintf("/?page=%d&per_page=1", page)
		req, err := env.Client.NewRequest(ctx, "GET", path, nil)
		if err != nil {
			return nil, err
		}
		resp, err := env.Client.Do(req, nil)
		if err != nil {
			return resp, err
		}
		expectedPage++
		return resp, err
	})

	if expectedPage != 4 {
		t.Errorf("expected to have walked through 3 pages, but walked through %d pages", expectedPage-1)
	}
}

func TestClientDo(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	callCount := 0
	env.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		switch callCount {
		case 1:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code:    string(ErrorCodeConflict),
					Message: "conflict",
				},
			})
		case 2:
			fmt.Fprintln(w, "{}")
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}
	})

	ctx := context.Background()
	request, _ := env.Client.NewRequest(ctx, http.MethodGet, "/test", nil)
	_, err := env.Client.Do(request, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("unexpected callCount: %v", callCount)
	}
}

func TestClientDoPost(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	debugLog := new(bytes.Buffer)

	env.Client.debugWriter = debugLog
	callCount := 0
	env.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer token" {
			t.Errorf("unexpected auth header: %q, expected %q", auth, "Bearer token")
		}

		callCount++
		w.Header().Set("Content-Type", "application/json")
		var dat map[string]interface{}
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			t.Error(err)
		}
		if err := json.Unmarshal(body, &dat); err != nil {
			t.Error(err)
		}
		switch callCount {
		case 1:
			if dat["test"] != "abcd" {
				t.Errorf("unexpected payload: %v", dat)
			}
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code:    string(ErrorCodeConflict),
					Message: "conflict",
				},
			})
		case 2:
			if dat["test"] != "abcd" {
				t.Errorf("unexpected payload: %v", dat)
			}
			fmt.Fprintln(w, "{}")
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}
	})

	ctx := context.Background()
	request, _ := env.Client.NewRequest(ctx, http.MethodPost, "/test", strings.NewReader(`{"test": "abcd"}`))
	_, err := env.Client.Do(request, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("unexpected callCount: %v", callCount)
	}

	if strings.Contains(debugLog.String(), "token") {
		t.Errorf("debug log did contain token, although it shouldn't")
	}
}

func TestBuildUserAgent(t *testing.T) {
	testCases := []struct {
		name               string
		applicationName    string
		applicationVersion string
		userAgent          string
	}{
		{"with application name and version", "test", "1.0", "test/1.0 " + UserAgent},
		{"with application name but no version", "test", "", "test " + UserAgent},
		{"without application name and version", "", "", UserAgent},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client := NewClient(WithApplication(testCase.applicationName, testCase.applicationVersion))
			if client.userAgent != testCase.userAgent {
				t.Errorf("unexpected user agent: %v", client.userAgent)
			}
		})
	}
}

func TestExponentialBackoff(t *testing.T) {
	t.Run("without jitter", func(t *testing.T) {
		backoffFunc := ExponentialBackoffWithOpts(ExponentialBackoffOpts{
			Base:       time.Second,
			Multiplier: 2,
			Cap:        32 * time.Second,
		})

		sum := 0.0
		for i, expected := range []time.Duration{
			1 * time.Second,
			2 * time.Second,
			4 * time.Second,
			8 * time.Second,
			16 * time.Second,
			32 * time.Second,
			32 * time.Second,
			32 * time.Second,
		} {
			backoff := backoffFunc(i)
			require.Equal(t, backoff, expected)
			sum += backoff.Seconds()
		}
		require.Equal(t, 127.0, sum)
	})

	t.Run("with jitter", func(t *testing.T) {
		backoffFunc := ExponentialBackoffWithOpts(ExponentialBackoffOpts{
			Base:       time.Second,
			Multiplier: 2,
			Cap:        32 * time.Second,
			Jitter:     true,
		})

		for i, expected := range []time.Duration{
			1 * time.Second,
			2 * time.Second,
			4 * time.Second,
			8 * time.Second,
			16 * time.Second,
			32 * time.Second,
			32 * time.Second,
			32 * time.Second,
		} {
			backoff := backoffFunc(i)
			assert.GreaterOrEqual(t, backoff, time.Second)
			assert.LessOrEqual(t, backoff, expected)
		}
	})
}
