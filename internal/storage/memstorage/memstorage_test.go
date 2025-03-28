package memstorage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
)

func CreateTestMemStorage() (storage *MemStorage) {
	storage = New()

	testCounter := metric.New("counter")
	testCounter.ID = "testCounter"
	*testCounter.Delta = 1

	_, _ = storage.UpdateMetric(testCounter)

	testGauge := metric.New("gauge")
	testGauge.ID = "testGauge"
	*testGauge.Value = 2.5

	_, _ = storage.UpdateMetric(testGauge)
	return
}

func TestNew(t *testing.T) {
	memStore := &MemStorage{
		metrics: make(map[string]metric.Metric),
	}
	got := New()
	require.Equal(t, got, memStore)
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	memStore := CreateTestMemStorage()
	want1 := [][2]string{
		[2]string{
			"testCounter",
			"1",
		},
		[2]string{
			"testGauge",
			"2.5",
		},
	}
	want2 := [][2]string{
		[2]string{
			"testGauge",
			"2.5",
		},
		[2]string{
			"testCounter",
			"1",
		},
	}
	want := [][][2]string{
		want1,
		want2,
	}

	got, _ := memStore.GetAllMetrics()
	assert.Contains(t, want, got)
}

func TestMemStorage_GetMetricValue(t *testing.T) {
	memStore := CreateTestMemStorage()

	wantCounter := metric.New("counter")
	wantCounter.ID = "testCounter"
	*wantCounter.Delta = 1

	wantGauge := metric.New("gauge")
	wantGauge.ID = "testGauge"
	*wantGauge.Value = 2.5

	t.Run("success get counter", func(t *testing.T) {
		got, err := memStore.GetMetricValue("testCounter", "counter")
		require.NoError(t, err)
		require.Equal(t, wantCounter, got)
	})
	t.Run("success get gauge", func(t *testing.T) {
		got, err := memStore.GetMetricValue("testGauge", "gauge")
		require.NoError(t, err)
		require.Equal(t, wantGauge, got)
	})
	t.Run("error wrong type counter", func(t *testing.T) {
		_, err := memStore.GetMetricValue("testCounter", "gauge")
		require.ErrorIs(t, err, metric.ErrInvalidType)
	})
	t.Run("error wrong type gauge", func(t *testing.T) {
		_, err := memStore.GetMetricValue("testGauge", "counter")
		require.ErrorIs(t, err, metric.ErrInvalidType)
	})
	t.Run("error unknown counter", func(t *testing.T) {
		_, err := memStore.GetMetricValue("testCounter2", "counter")
		require.ErrorIs(t, err, metric.ErrUnknownMetric)
	})
	t.Run("error unknown gauge", func(t *testing.T) {
		_, err := memStore.GetMetricValue("testGauge2", "gauge")
		require.ErrorIs(t, err, metric.ErrUnknownMetric)
	})
}

func TestMemStorage_UpdateMetric(t *testing.T) {
	store := CreateTestMemStorage()

	updateCounter := metric.New("counter")
	updateCounter.ID = "testCounter"
	*updateCounter.Delta = 3

	wantCounter := metric.New("counter")
	wantCounter.ID = "testCounter"
	*wantCounter.Delta = 4

	wantGauge := metric.New("gauge")
	wantGauge.ID = "testGauge"
	*wantGauge.Value = 5.5

	noName := metric.New("counter")

	noType := metric.Metric{
		ID: "noType",
	}

	noVal := metric.Metric{
		ID:    "noType",
		MType: "counter",
	}

	newCounter := metric.New("counter")
	newCounter.ID = "testCounterNew"
	*newCounter.Delta = 2

	newGauge := metric.New("gauge")
	newGauge.ID = "testGaugeNew"
	*newGauge.Value = 3.5

	invalidTypeGauge := metric.New("counter")
	invalidTypeGauge.ID = "testGauge"

	tests := []struct {
		name    string
		s       *MemStorage
		arg     metric.Metric
		want    metric.Metric
		wantErr bool
		err     error
	}{
		{
			name:    "success update counter",
			s:       store,
			arg:     updateCounter,
			want:    wantCounter,
			wantErr: false,
		},
		{
			name:    "success update gauge",
			s:       store,
			arg:     wantGauge,
			want:    wantGauge,
			wantErr: false,
		},
		{
			name:    "error no name metric",
			s:       store,
			arg:     noName,
			wantErr: true,
			err:     metric.ErrInvalidName,
		},
		{
			name:    "error no type metric",
			s:       store,
			arg:     noType,
			wantErr: true,
			err:     metric.ErrInvalidType,
		},
		{
			name:    "error no value metric",
			s:       store,
			arg:     noVal,
			wantErr: true,
			err:     metric.ErrParseValue,
		},
		{
			name:    "success create new counter metric",
			s:       store,
			arg:     newCounter,
			wantErr: false,
			want:    newCounter,
		},
		{
			name:    "success create new gauge metric",
			s:       store,
			arg:     newGauge,
			wantErr: false,
			want:    newGauge,
		},
		{
			name:    "error invalid type",
			s:       store,
			arg:     invalidTypeGauge,
			wantErr: true,
			err:     metric.ErrInvalidType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.UpdateMetric(tt.arg)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
				return
			}
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestMemStorage_UpdateManyMetrics(t *testing.T) {
	store := CreateTestMemStorage()

	updateCounter := metric.New("counter")
	updateCounter.ID = "testCounter"
	*updateCounter.Delta = 2

	wantCounter := metric.New("counter")
	wantCounter.ID = "testCounter"
	*wantCounter.Delta = 3

	updateGauge := metric.New("gauge")
	updateGauge.ID = "testGauge"
	*updateGauge.Value = 5.5

	newGauge := metric.New("gauge")
	newGauge.ID = "testGaugeNew"
	*newGauge.Value = 6.4

	forUpdate := []metric.Metric{
		updateCounter,
		updateGauge,
		newGauge,
	}

	t.Run("success", func(t *testing.T) {
		err := store.UpdateManyMetrics(context.TODO(), forUpdate)
		require.NoError(t, err)

		counter, err := store.GetMetricValue("testCounter", "counter")
		require.NoError(t, err)
		require.Equal(t, wantCounter, counter)

		gauge, err := store.GetMetricValue("testGauge", "gauge")
		require.NoError(t, err)
		require.Equal(t, updateGauge, gauge)

		gaugeNew, err := store.GetMetricValue("testGaugeNew", "gauge")
		require.NoError(t, err)
		require.Equal(t, newGauge, gaugeNew)
	})

	errMetric := metric.Metric{
		ID:    "noValMetric",
		MType: "counter",
	}
	forUpdate = append(forUpdate, errMetric)

	t.Run("error update one of metric", func(t *testing.T) {
		err := store.UpdateManyMetrics(context.TODO(), forUpdate)
		require.Error(t, err)
	})
}

func TestMemStorage_CreateMemento(t *testing.T) {
	s := CreateTestMemStorage()
	mem := s.CreateMemento()
	require.Equal(t, s.metrics, mem.Metrics)
}

func TestMemento_GetSavedState(t *testing.T) {
	s := CreateTestMemStorage()
	mem := s.CreateMemento()
	cur := mem.GetSavedState()
	require.Equal(t, s.metrics, cur)
}

func TestMemStorage_RestoreMemento(t *testing.T) {
	s := CreateTestMemStorage()
	mem := s.CreateMemento()

	copyS := CreateTestMemStorage()

	newmetric := metric.New("counter")
	newmetric.ID = "newCounter"
	*newmetric.Delta = 5

	_, _ = copyS.UpdateMetric(newmetric)

	s.RestoreMemento(mem)

	require.Equal(t, s.metrics, mem.Metrics)
}

func TestMemStorage_Close(t *testing.T) {
	s := CreateTestMemStorage()
	err := s.Close()
	require.NoError(t, err)
}

func TestMemStorage_Ping(t *testing.T) {
	s := CreateTestMemStorage()
	err := s.Ping()
	require.NoError(t, err)
}
