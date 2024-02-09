package instrumentation

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func Test_preparePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			"simple test",
			"/v1/volumes/123456",
			"/volumes/",
		},
		{
			"simple test",
			"/v1/volumes/123456/actions/attach",
			"/volumes/actions/attach",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := preparePathForLabel(tt.path); got != tt.want {
				t.Errorf("preparePathForLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultipleInstrumentedClients(t *testing.T) {
	reg := prometheus.NewRegistry()

	t.Run("should not panic", func(_ *testing.T) {
		// Following code should run without panicking
		New("test", reg).InstrumentedRoundTripper()
		New("test", reg).InstrumentedRoundTripper()
	})
}
