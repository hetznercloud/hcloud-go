package hcloud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterPages(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		result, err := iterPages(func(page int) ([]*int, *Response, error) {
			if page < 4 {
				return []*int{new(page)}, &Response{Meta: Meta{Pagination: &Pagination{NextPage: page + 1}}}, nil
			}
			return []*int{new(page)}, &Response{}, nil
		})
		require.NoError(t, err)
		require.Equal(t, []*int{new(1), new(2), new(3), new(4)}, result)
	})

	t.Run("failed", func(t *testing.T) {
		result, err := iterPages(func(page int) ([]*int, *Response, error) {
			if page < 4 {
				return []*int{new(page)}, &Response{Meta: Meta{Pagination: &Pagination{NextPage: page + 1}}}, nil
			}
			return nil, &Response{}, fmt.Errorf("failure")
		})
		require.EqualError(t, err, "failure")
		require.Nil(t, result)
	})
}
