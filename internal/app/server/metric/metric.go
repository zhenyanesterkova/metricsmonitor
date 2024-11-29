package metric

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Metric struct {
	MetricGauge
	MetricCounter
	ID    string `json:"id"`
	MType string `json:"type"`
}

func (m *Metric) GetType() string {
	return m.MType
}

func (m *Metric) String() string {
	switch m.MType {
	case TypeGauge:
		return m.MetricGauge.String()
	case TypeCounter:
		return m.MetricCounter.String()
	default:
		return ""
	}
}

func New(typeMetric string) Metric {
	switch typeMetric {
	case TypeGauge:
		return Metric{
			MType:       TypeGauge,
			MetricGauge: NewMetricGauge(),
		}
	case TypeCounter:
		return Metric{
			MType:         TypeCounter,
			MetricCounter: NewMetricCounter(),
		}
	default:
		return Metric{
			MetricGauge:   NewMetricGauge(),
			MetricCounter: NewMetricCounter(),
		}
	}
}
