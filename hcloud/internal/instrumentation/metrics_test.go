package instrumentation

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMultipleInstrumentedClients(t *testing.T) {
	reg := prometheus.NewRegistry()

	t.Run("should not panic", func(_ *testing.T) {
		// Following code should run without panicking
		New("test", reg).InstrumentedRoundTripper()
		New("test", reg).InstrumentedRoundTripper()
	})
}

func TestPreparePathForLabel(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{
			"/v1/volumes/123456",
			"/volumes/-",
		},
		{
			"/v1/volumes/123456/actions/attach",
			"/volumes/-/actions/attach",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.want, preparePathForLabel(tt.path))
		})
	}
}
