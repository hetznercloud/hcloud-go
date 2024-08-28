package actionutil

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestAppendNext(t *testing.T) {
	action := &hcloud.Action{ID: 1}
	nextActions := []*hcloud.Action{{ID: 2}, {ID: 3}}

	actions := AppendNext(action, nextActions)

	assert.Equal(t, []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}}, actions)
}

func TestAllForResource(t *testing.T) {
	ctx := context.Background()

	server := httptest.NewServer(mockutil.Handler(t, []mockutil.Request{
		{
			Method: "GET", Path: "/firewalls/actions?page=1&status=running",
			Status: 200,
			JSON: schema.ActionListResponse{
				Actions: []schema.Action{
					{Resources: []schema.ActionResourceReference{{Type: "server", ID: 8}}},
					{Resources: []schema.ActionResourceReference{{Type: "server", ID: 8}, {Type: "firewall", ID: 10}}},
					{Resources: []schema.ActionResourceReference{{Type: "server", ID: 10}}},
					{Resources: []schema.ActionResourceReference{{Type: "firewall", ID: 8}}},
				},
			},
		},
	}))
	client := hcloud.NewClient(hcloud.WithEndpoint(server.URL))

	actions, err := AllForResource(ctx,
		client.Firewall.Action,
		hcloud.ActionListOpts{Status: []hcloud.ActionStatus{hcloud.ActionStatusRunning}},
		hcloud.ActionResourceTypeServer, 8,
	)
	assert.NoError(t, err)
	assert.Len(t, actions, 2)
}
