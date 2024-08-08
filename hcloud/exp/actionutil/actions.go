package actionutil

import (
	"context"
	"slices"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// AppendNext return the action and the next actions in a new slice.
func AppendNext(action *hcloud.Action, nextActions []*hcloud.Action) []*hcloud.Action {
	all := make([]*hcloud.Action, 0, 1+len(nextActions))
	all = append(all, action)
	all = append(all, nextActions...)
	return all
}

func FindForResource(action hcloud.Action, kind hcloud.ActionResourceType, id int64) *hcloud.ActionResource {
	i := slices.IndexFunc(action.Resources, func(r *hcloud.ActionResource) bool {
		return r.Type == kind && r.ID == id
	})

	if i > -1 {
		return action.Resources[i]
	}
	return nil
}

func RunningForResource(
	client *hcloud.Client,
	ctx context.Context,
	kind hcloud.ActionResourceType,
	id int64,
) ([]*hcloud.Action, error) {
	actions, err := client.Server.Action.All(ctx,
		hcloud.ActionListOpts{
			Status: []hcloud.ActionStatus{hcloud.ActionStatusRunning},
		},
	)
	if err != nil {
		return nil, err
	}

	actions = slices.Clip(slices.DeleteFunc(actions, func(i *hcloud.Action) bool {
		return i.Resource.Type == kind && i.Resource.ID == id
	}))

	return actions, nil
}
