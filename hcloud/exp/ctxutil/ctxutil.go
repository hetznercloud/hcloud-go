package ctxutil

import (
	"context"
	"strings"
)

// key is an unexported type to prevents collisions with keys defined in other packages.
type key struct{}

// opPathKey is the key for operation path in Contexts.
var opPathKey = key{}

func SetOpPath(ctx context.Context, path string) context.Context {
	path, _, _ = strings.Cut(path, "?")
	path = strings.ReplaceAll(path, "%d", "-")
	path = strings.ReplaceAll(path, "%s", "-")

	return context.WithValue(ctx, opPathKey, path)
}

func OpPath(ctx context.Context) string {
	result := ctx.Value(opPathKey).(string)
	return result
}
