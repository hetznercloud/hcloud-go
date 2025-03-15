package ctxutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpPath(t *testing.T) {
	for _, tt := range []struct {
		path string
		want string
	}{
		{
			path: "/resource/%s/nested/%s/%d",
			want: "/resource/-/nested/-/-",
		},
		{
			path: "/certificates/%d",
			want: "/certificates/-",
		},
		{
			path: "/servers/%d/metrics?%s",
			want: "/servers/-/metrics",
		},
	} {
		t.Run("", func(t *testing.T) {
			ctx := context.Background()

			require.Equal(t, tt.want, OpPath(SetOpPath(ctx, tt.path)))
		})
	}
}
