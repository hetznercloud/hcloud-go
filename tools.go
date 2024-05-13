//go:build tools
// +build tools

package tools

import (
	_ "github.com/jmattheis/goverter/cmd/goverter"
	_ "github.com/vburenin/ifacemaker"
	_ "go.uber.org/mock/mockgen"
)
