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

func AllForResource(
	ctx context.Context,
	actionClient *hcloud.ResourceActionClient,
	opts hcloud.ActionListOpts,
	resourceType hcloud.ActionResourceType,
	resourceID int64,
) ([]*hcloud.Action, error) {
	actions, err := actionClient.All(ctx, opts)
	if err != nil {
		return nil, err
	}

	actions = slices.Clip(slices.DeleteFunc(actions, func(action *hcloud.Action) bool {
		return !slices.ContainsFunc(action.Resources, func(resource *hcloud.ActionResource) bool {
			return resource.Type == resourceType && resource.ID == resourceID
		})
	}))

	return actions, nil
}
