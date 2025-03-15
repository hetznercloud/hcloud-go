package instrumentation

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMultipleInstrumentedClients(t *testing.T) {
	reg := prometheus.NewRegistry()

	t.Run("should not panic", func(_ *testing.T) {
		// Following code should run without panicking
		New("test", reg).InstrumentedRoundTripper()
		New("test", reg).InstrumentedRoundTripper()
	})
}
