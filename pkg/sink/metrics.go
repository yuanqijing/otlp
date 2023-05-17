package sink

import (
	"context"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"sync"
)

var _ consumer.Metrics = (*MetricsSink)(nil)

// MetricsSink is a consumer.Metrics that acts like a sink that
// stores all metrics and allows querying them for testing.
type MetricsSink struct {
	mu             sync.Mutex
	metrics        []pmetric.Metrics
	dataPointCount int
}

func (sme *MetricsSink) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeMetrics stores metrics to this sink.
func (sme *MetricsSink) ConsumeMetrics(_ context.Context, md pmetric.Metrics) error {
	sme.mu.Lock()
	defer sme.mu.Unlock()

	sme.metrics = append(sme.metrics, md)
	sme.dataPointCount += md.DataPointCount()

	return nil
}

// AllMetrics returns the metrics stored by this sink since last Reset.
func (sme *MetricsSink) AllMetrics() []pmetric.Metrics {
	sme.mu.Lock()
	defer sme.mu.Unlock()

	copyMetrics := make([]pmetric.Metrics, len(sme.metrics))
	copy(copyMetrics, sme.metrics)
	return copyMetrics
}

// DataPointCount returns the number of metrics stored by this sink since last Reset.
func (sme *MetricsSink) DataPointCount() int {
	sme.mu.Lock()
	defer sme.mu.Unlock()
	return sme.dataPointCount
}

// Reset deletes any stored data.
func (sme *MetricsSink) Reset() {
	sme.mu.Lock()
	defer sme.mu.Unlock()

	sme.metrics = nil
	sme.dataPointCount = 0
}
