package sliceutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransform(t *testing.T) {
	require.Equal(t,
		[]string{"1", "2"},
		Transform(
			[]int{1, 2},
			func(e int) string { return fmt.Sprint(e) },
		),
	)
}
