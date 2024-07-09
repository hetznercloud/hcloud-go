package retryutils

import (
	"context"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutils"
)

func FakeBackoff(counter *int64) hcloud.BackoffFunc {
	return func(_ int) time.Duration {
		atomic.AddInt64(counter, 1)
		return 0
	}
}

func TestRetry(t *testing.T) {
	sshKey := &hcloud.SSHKey{ID: 123}

	updateResponse := mockutils.Request{
		Method: "PUT", Path: "/ssh_keys/123",
		Status:  200,
		JSONRaw: `{ "ssh_key": { "id": 123 }}`,
	}
	updateLockedResponse := mockutils.Request{
		Method: "PUT", Path: "/ssh_keys/123",
		Status:  423,
		JSONRaw: `{ "error": { "code": "locked", "message": "Resource is locked" }}`,
	}
	deleteResponse := mockutils.Request{
		Method: "DELETE", Path: "/ssh_keys/123",
		Status: 204,
	}
	deleteLockedResponse := mockutils.Request{
		Method: "DELETE", Path: "/ssh_keys/123",
		Status:  423,
		JSONRaw: `{ "error": { "code": "locked", "message": "Resource is locked"  }}`,
	}

	testCases := []struct {
		name     string
		requests []mockutils.Request
		run      func(t *testing.T, client *hcloud.Client)
	}{
		{
			name: "happy with 3 return values",
			requests: []mockutils.Request{
				updateLockedResponse,
				updateLockedResponse,
				updateResponse,
			},
			run: func(t *testing.T, client *hcloud.Client) {
				var retryCount int64
				retryOpts := Opts{
					Backoff:    FakeBackoff(&retryCount),
					MaxRetries: 3,
					Policy: func(err error) bool {
						return hcloud.IsError(err, hcloud.ErrorCodeLocked)
					},
				}

				result, resp, err := Retry(retryOpts, func() (*hcloud.SSHKey, *hcloud.Response, error) {
					return client.SSHKey.Update(context.Background(), sshKey, hcloud.SSHKeyUpdateOpts{})
				})
				require.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, int64(123), result.ID)

				assert.Equal(t, int64(2), retryCount)
			},
		},
		{
			name: "happy with 2 return values",
			requests: []mockutils.Request{
				deleteLockedResponse,
				deleteLockedResponse,
				deleteResponse,
			},
			run: func(t *testing.T, client *hcloud.Client) {
				var retryCount int64
				retryOpts := Opts{
					Backoff:    FakeBackoff(&retryCount),
					MaxRetries: 3,
					Policy: func(err error) bool {
						return hcloud.IsError(err, hcloud.ErrorCodeLocked)
					},
				}

				resp, err := RetryNoResult(retryOpts, func() (*hcloud.Response, error) {
					return client.SSHKey.Delete(context.Background(), sshKey)
				})
				require.NoError(t, err)
				assert.NotNil(t, resp)

				assert.Equal(t, int64(2), retryCount)
			},
		},
		{
			name: "fail with locked error",
			requests: []mockutils.Request{
				deleteLockedResponse,
				deleteLockedResponse,
				deleteLockedResponse,
				deleteLockedResponse,
			},
			run: func(t *testing.T, client *hcloud.Client) {
				var retryCount int64
				retryOpts := Opts{
					Backoff:    FakeBackoff(&retryCount),
					MaxRetries: 3,
					Policy: func(err error) bool {
						return hcloud.IsError(err, hcloud.ErrorCodeLocked)
					},
				}

				resp, err := RetryNoResult(retryOpts, func() (*hcloud.Response, error) {
					return client.SSHKey.Delete(context.Background(), sshKey)
				})
				require.EqualError(t, err, "Resource is locked (locked)")
				assert.NotNil(t, resp)

				assert.Equal(t, int64(3), retryCount)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(mockutils.Handler(t, testCase.requests))
			defer server.Close()

			client := hcloud.NewClient(hcloud.WithEndpoint(server.URL))

			testCase.run(t, client)
		})
	}
}
