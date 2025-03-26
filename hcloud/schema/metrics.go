package schema

import (
	"time"

	"github.com/prometheus/common/model"
)

type Metrics struct {
	Start      time.Time                   `json:"start"`
	End        time.Time                   `json:"end"`
	Step       float64                     `json:"step"`
	TimeSeries map[string]MetricsTimeSerie `json:"time_series"`
}

// MetricsTimeSerie is a partial implementation of [model.SampleStream].
type MetricsTimeSerie struct {
	Values []MetricsTimeSerieValue `json:"values"`
}

type MetricsTimeSerieValue model.SamplePair
