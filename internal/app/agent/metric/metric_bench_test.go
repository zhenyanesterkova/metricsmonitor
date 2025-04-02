package metric

import "testing"

func BenchmarkMetric(b *testing.B) {
	metricGauge := &Metric{
		Value: new(float64),
		ID:    "gauge",
		MType: "gauge",
	}
	*metricGauge.Value = 3.5

	metricCounter := &Metric{
		Delta: new(int64),
		ID:    "counter",
		MType: "counter",
	}
	*metricCounter.Delta = 3

	b.Run("Counter: String()", func(b *testing.B) {
		for range b.N {
			_ = metricCounter.String()
		}
	})

	b.Run("Gauge: String()", func(b *testing.B) {
		for range b.N {
			_ = metricGauge.String()
		}
	})

	b.Run("Counter: StringValue()", func(b *testing.B) {
		for range b.N {
			_ = metricCounter.StringValue()
		}
	})

	b.Run("Gauge: StringValue()", func(b *testing.B) {
		for range b.N {
			_ = metricGauge.StringValue()
		}
	})

	b.Run("Gauge: updateGauge()", func(b *testing.B) {
		for range b.N {
			metricGauge.updateGauge(3.5)
		}
	})

	b.Run("Counter: updateCounter()", func(b *testing.B) {
		for range b.N {
			metricCounter.updateCounter()
		}
	})
}
