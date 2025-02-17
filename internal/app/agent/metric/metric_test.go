package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateGauge(t *testing.T) {
	metrica := Metric{
		ID:    "test",
		MType: "gauge",
	}
	temp := 3.5
	expectMetrica := Metric{
		ID:    "test",
		MType: "gauge",
		Value: &temp,
	}
	t.Run("Test #1", func(t *testing.T) {
		metrica.updateGauge(3.5)
		assert.Equal(t, expectMetrica, metrica)
	})
}

func Test_updateCounter(t *testing.T) {
	metrica := Metric{
		ID:    "test",
		MType: "counter",
	}
	temp := int64(1)
	expectMetrica := Metric{
		ID:    "test",
		MType: "counter",
		Delta: &temp,
	}
	t.Run("Test #1", func(t *testing.T) {
		metrica.updateCounter()
		assert.Equal(t, expectMetrica, metrica)
	})
}
