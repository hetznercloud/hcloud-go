package errutil

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// LogValue returns a [slog.Value] for a [hcloud.Error].
func LogValue(err error) slog.Value {
	var herr *hcloud.Error
	if !errors.As(err, &herr) {
		return slog.AnyValue(err)
	}

	attrs := []slog.Attr{
		slog.String("msg", herr.Error()),
	}

	if herr.Details != nil {
		switch details := herr.Details.(type) {
		case hcloud.ErrorDetailsInvalidInput:
			attrs = append(attrs,
				slog.String("details", fmt.Sprintf("%v", details.Fields)),
			)
		case hcloud.ErrorDetailsDeprecatedAPIEndpoint:
			attrs = append(attrs,
				slog.String("details", fmt.Sprintf("%v", details)),
			)
		default:
			attrs = append(attrs,
				slog.String("details", fmt.Sprintf("%v", details)),
			)
		}
	}

	return slog.GroupValue(attrs...)
}
