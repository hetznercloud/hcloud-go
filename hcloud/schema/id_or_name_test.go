package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIDOrNameMarshall(t *testing.T) {
	t.Run("id", func(t *testing.T) {
		i := IDOrName{ID: 1}

		got, err := i.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `1`, string(got))
	})

	t.Run("name", func(t *testing.T) {
		i := IDOrName{Name: "name"}

		got, err := i.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `"name"`, string(got))
	})

	t.Run("id and name", func(t *testing.T) {
		i := IDOrName{ID: 1, Name: "name"}

		got, err := i.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `1`, string(got))
	})

	t.Run("none", func(t *testing.T) {
		i := IDOrName{}

		_, err := i.MarshalJSON()
		require.EqualError(t, err, "json: unsupported value: id or name must not be zero values")
	})
}

func TestIDOrNameUnMarshall(t *testing.T) {
	t.Run("id", func(t *testing.T) {
		i := IDOrName{}

		err := i.UnmarshalJSON([]byte(`1`))
		require.NoError(t, err)
		require.Equal(t, IDOrName{ID: 1}, i)
	})
	t.Run("name", func(t *testing.T) {
		i := IDOrName{}

		err := i.UnmarshalJSON([]byte(`"name"`))
		require.NoError(t, err)
		require.Equal(t, IDOrName{Name: "name"}, i)
	})
	t.Run("id string", func(t *testing.T) {
		i := IDOrName{}

		err := i.UnmarshalJSON([]byte(`"1"`))
		require.NoError(t, err)
		require.Equal(t, IDOrName{ID: 1}, i)
	})
	t.Run("id float", func(t *testing.T) {
		i := IDOrName{}

		err := i.UnmarshalJSON([]byte(`1.0`))
		require.EqualError(t, err, "json: cannot unmarshal 1.0 into Go value of type schema.IDOrName")
	})
	t.Run("null", func(t *testing.T) {
		i := IDOrName{}

		err := i.UnmarshalJSON([]byte(`null`))
		require.EqualError(t, err, "json: cannot unmarshal null into Go value of type schema.IDOrName")
	})
}
