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
				return []*int{Ptr(page)}, &Response{Meta: Meta{Pagination: &Pagination{NextPage: page + 1}}}, nil
			}
			return []*int{Ptr(page)}, &Response{}, nil
		})
		require.NoError(t, err)
		require.Equal(t, []*int{Ptr(1), Ptr(2), Ptr(3), Ptr(4)}, result)
	})

	t.Run("failed", func(t *testing.T) {
		result, err := iterPages(func(page int) ([]*int, *Response, error) {
			if page < 4 {
				return []*int{Ptr(page)}, &Response{Meta: Meta{Pagination: &Pagination{NextPage: page + 1}}}, nil
			}
			return nil, &Response{}, fmt.Errorf("failure")
		})
		require.EqualError(t, err, "failure")
		require.Nil(t, result)
	})
}
