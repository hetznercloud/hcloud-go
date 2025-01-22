package randutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomID(t *testing.T) {
	found1 := GenerateID()
	found2 := GenerateID()

	assert.Len(t, found1, 8)
	assert.Len(t, found2, 8)
	assert.NotEqual(t, found1, found2)
}
