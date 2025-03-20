package schema

import (
	"encoding/json"
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

	t.Run("null", func(t *testing.T) {
		i := IDOrName{}

		got, err := i.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `null`, string(got))
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

func TestIDOrName(t *testing.T) {
	// Make sure the behavior does not change from the use of an interface{}.
	type FakeRequest struct {
		Old any      `json:"old"`
		New IDOrName `json:"new"`
	}

	t.Run("null", func(t *testing.T) {
		o := FakeRequest{}
		body, err := json.Marshal(o)
		require.NoError(t, err)
		require.JSONEq(t, `{"old":null,"new":null}`, string(body))
	})
	t.Run("id", func(t *testing.T) {
		o := FakeRequest{Old: int64(1), New: IDOrName{ID: 1}}
		body, err := json.Marshal(o)
		require.NoError(t, err)
		require.JSONEq(t, `{"old":1,"new":1}`, string(body))
	})
	t.Run("name", func(t *testing.T) {
		o := FakeRequest{Old: "name", New: IDOrName{Name: "name"}}
		body, err := json.Marshal(o)
		require.NoError(t, err)
		require.JSONEq(t, `{"old":"name","new":"name"}`, string(body))
	})
}
