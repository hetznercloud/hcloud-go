package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAppendNextActions(t *testing.T) {
	action := &hcloud.Action{ID: 1}
	nextActions := []*hcloud.Action{{ID: 2}, {ID: 3}}

	actions := AppendNextActions(action, nextActions)

	assert.Equal(t, []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}}, actions)
}
