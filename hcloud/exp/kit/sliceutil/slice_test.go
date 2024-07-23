package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatches(t *testing.T) {
	all := []int{1, 2, 3, 4, 5}
	batches := Batches(all, 2)

	assert.Len(t, batches, 3)

	assert.Equal(t, []int{1, 2}, batches[0])
	assert.Equal(t, 2, cap(batches[0]))

	assert.Equal(t, []int{3, 4}, batches[1])
	assert.Equal(t, 2, cap(batches[1]))

	assert.Equal(t, []int{5}, batches[2])
	assert.Equal(t, 1, cap(batches[2]))
}
