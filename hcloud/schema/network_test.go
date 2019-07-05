package schema

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNetworkUpdateRequest(t *testing.T) {
	var (
		oneLabel    = map[string]string{"foo": "bar"}
		nilLabels   map[string]string
		emptyLabels = map[string]string{}
	)

	testCases := []struct {
		name string
		in   NetworkUpdateRequest
		out  []byte
	}{
		{
			name: "no labels",
			in:   NetworkUpdateRequest{},
			out:  []byte(`{}`),
		},
		{
			name: "one label",
			in:   NetworkUpdateRequest{Labels: &oneLabel},
			out:  []byte(`{"labels":{"foo":"bar"}}`),
		},
		{
			name: "nil labels",
			in:   NetworkUpdateRequest{Labels: &nilLabels},
			out:  []byte(`{"labels":null}`),
		},
		{
			name: "empty labels",
			in:   NetworkUpdateRequest{Labels: &emptyLabels},
			out:  []byte(`{"labels":{}}`),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := json.Marshal(testCase.in)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(data, testCase.out) {
				t.Fatalf("output %s does not match %s", data, testCase.out)
			}
		})
	}
}
