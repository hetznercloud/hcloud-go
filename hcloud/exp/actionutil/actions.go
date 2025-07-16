package actionutil

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

// AppendNext return the action and the next actions in a new slice.
//
// Experimental: Breaking changes may occur within minor releases.
func AppendNext(action *hcloud.Action, nextActions []*hcloud.Action) []*hcloud.Action {
	all := make([]*hcloud.Action, 0, 1+len(nextActions))
	all = append(all, action)
	all = append(all, nextActions...)
	return all
}
