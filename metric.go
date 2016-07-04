package ddleash

import (
	"encoding/json"
	"fmt"
)

type MetricType int

const (
	Gauge MetricType = iota
	Rate
)

type Metric struct {
	Name        string
	Description string
	Type        MetricType `json:"metric_type"`
	Interval    uint
}

func (metricType *MetricType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	got, ok := map[string]MetricType{"gauge": Gauge, "rate": Rate}[s]
	if !ok {
		return fmt.Errorf("Unknwon metric type %q", s)
	}

	*metricType = got
	return nil
}
