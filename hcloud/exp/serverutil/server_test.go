package serverutil

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestEnsurePowerOff(t *testing.T) {
	for _, testCase := range []struct {
		Name         string
		WantRequests []mockutil.Request
		Run          func(client *hcloud.Client)
	}{
		{
			Name: "succeed",
			WantRequests: []mockutil.Request{
				{
					Method: "POST", Path: "/servers/1/actions/poweroff",
					Status: 200,
					JSONRaw: `{
						"action": { "id": 1509772237, "status": "running", "progress": 0 }
					}`,
				},
				{
					Method: "GET", Path: "/actions?id=1509772237&page=1&sort=status&sort=id",
					Status: 200,
					JSONRaw: `{
						"actions": [
							{ "id": 1509772237, "status": "success", "progress": 100 }
						]
					}`,
				},
			},
			Run: func(client *hcloud.Client) {
				ctx := context.Background()
				err := EnsurePowerOff(ctx, client, &hcloud.Server{ID: 1, Status: hcloud.ServerStatusRunning})
				assert.NoError(t, err)
			},
		},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			server := httptest.NewServer(mockutil.Handler(t, testCase.WantRequests))
			defer server.Close()

			client := hcloud.NewClient(
				hcloud.WithEndpoint(server.URL),
				hcloud.WithPollOpts(hcloud.PollOpts{BackoffFunc: hcloud.ConstantBackoff(0)}),
			)

			testCase.Run(client)
		})
	}
}
