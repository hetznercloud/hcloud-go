package labelutils

import "testing"

func TestSelector(t *testing.T) {
	tests := []struct {
		name             string
		labels           map[string]string
		expectedSelector string
	}{
		{
			name:             "empty selector",
			labels:           map[string]string{},
			expectedSelector: "",
		},
		{
			name:             "single label",
			labels:           map[string]string{"foo": "bar"},
			expectedSelector: "foo=bar",
		},
		{
			name:             "multiple labels",
			labels:           map[string]string{"foo": "bar", "foz": "baz"},
			expectedSelector: "foo=bar,foz=baz",
		},
		{
			name:             "nil map",
			labels:           nil,
			expectedSelector: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Selector(tt.labels); got != tt.expectedSelector {
				t.Errorf("Selector() = %v, want %v", got, tt.expectedSelector)
			}
		})
	}
}
