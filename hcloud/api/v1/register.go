package v1

import (
	"github.com/hetznercloud/hcloud-go/hcloud/api"
	"github.com/hetznercloud/hcloud-go/hcloud/runtime"
	"github.com/hetznercloud/hcloud-go/hcloud/runtime/conversion"
)

var (
	// Scheme of the v1 API
	Scheme *runtime.Scheme
)

func init() {
	Scheme = runtime.NewScheme()
	Scheme.AddConversionFuncs(
		func(in *Action, out *api.Action, s conversion.Scope) error {
			if in.Error != nil {
				out.ErrorCode = in.Error.Code
				out.ErrorMessage = in.Error.Message
			}
			return s.DefaultConvert(in, out)
		},
	)
}
