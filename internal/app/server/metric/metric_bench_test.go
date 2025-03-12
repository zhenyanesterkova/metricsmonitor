package metric

import "testing"

func BenchmarkMetric(b *testing.B) {
	b.Run("Metric: NewMetric()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New(TypeCounter)
		}
	})

	metrica := New(TypeCounter)

	b.Run("Metric: GetType()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = metrica.GetType()
		}
	})

	b.Run("Metric: String()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = metrica.String()
		}
	})

	b.Run("Gauge: NewMetricGauge()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewMetricGauge()
		}
	})

	metricGauge := NewMetricGauge()

	b.Run("Gauge: SetValue()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metricGauge.SetValue(3.5)
		}
	})

	b.Run("Gauge: String()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = metricGauge.String()
		}
	})

	b.Run("Counter: NewMetricCounter()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewMetricCounter()
		}
	})

	metricCounter := NewMetricCounter()

	b.Run("Counter: SetValue()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			metricCounter.SetValue(3)
		}
	})

	b.Run("Counter: String()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = metricCounter.String()
		}
	})
}
