package transaction

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/kit/mocked"
	"github.com/stretchr/testify/assert"
)

func TestWithServerOff(t *testing.T) {
	server := &hcloud.Server{ID: 1}

	for _, testCase := range []struct {
		Name         string
		WantRequests []mocked.Request
		Run          func(client *hcloud.Client)
	}{
		{
			Name: "succeed",
			WantRequests: []mocked.Request{
				{Method: "POST", Path: "/servers/1/actions/poweroff",
					Status: 200,
					Body: `{
						"actions": [
							{ "id": 1509772237, "status": "running", "progress": 0 }
						]
					}`},
				{Method: "POST", Path: "/servers/1/actions/poweron",
					Status: 200,
					Body: `{
						"actions": [
							{ "id": 1509772237, "status": "running", "progress": 0 }
						]
					}`},
			},
			Run: func(client *hcloud.Client) {
				ctx := context.Background()

				nextCalled := false
				err := WithServerPowerOff(ctx, client, server, func() error {
					nextCalled = true
					return nil
				})

				assert.True(t, nextCalled)
				assert.NoError(t, err)
			},
		},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			server := httptest.NewServer(mocked.Handler(t, testCase.WantRequests))
			defer server.Close()

			client := hcloud.NewClient(
				hcloud.WithEndpoint(server.URL),
				hcloud.WithBackoffFunc(func(_ int) time.Duration { return 0 }),
				hcloud.WithPollBackoffFunc(func(_ int) time.Duration { return 0 }),
			)

			testCase.Run(client)
		})
	}
}
