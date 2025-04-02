package metric

import "testing"

func BenchmarkMetric(b *testing.B) {
	b.Run("Metric: NewMetric()", func(b *testing.B) {
		for range b.N {
			_ = New(TypeCounter)
		}
	})

	metrica := New(TypeCounter)

	b.Run("Metric: GetType()", func(b *testing.B) {
		for range b.N {
			_ = metrica.GetType()
		}
	})

	b.Run("Metric: String()", func(b *testing.B) {
		for range b.N {
			_ = metrica.String()
		}
	})

	b.Run("Gauge: NewMetricGauge()", func(b *testing.B) {
		for range b.N {
			_ = NewMetricGauge()
		}
	})

	metricGauge := NewMetricGauge()

	b.Run("Gauge: SetValue()", func(b *testing.B) {
		for range b.N {
			metricGauge.SetValue(3.5)
		}
	})

	b.Run("Gauge: String()", func(b *testing.B) {
		for range b.N {
			_ = metricGauge.String()
		}
	})

	b.Run("Counter: NewMetricCounter()", func(b *testing.B) {
		for range b.N {
			_ = NewMetricCounter()
		}
	})

	metricCounter := NewMetricCounter()

	b.Run("Counter: SetValue()", func(b *testing.B) {
		for range b.N {
			metricCounter.SetValue(3)
		}
	})

	b.Run("Counter: String()", func(b *testing.B) {
		for range b.N {
			_ = metricCounter.String()
		}
	})
}
