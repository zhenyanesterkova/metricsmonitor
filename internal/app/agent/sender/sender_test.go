package sender

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/agent/metric"
)

func ServerMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestSender_SendQueryUpdateMetrics(t *testing.T) {
	hand := http.HandlerFunc(ServerMock)
	srv := httptest.NewServer(hand)

	metricBuff := metric.NewMetricBuf()
	metricBuff.UpdateMetrics()
	err := metricBuff.UpdateGopsutilMetrics()
	assert.NoError(t, err)

	key := "testong"
	t.Run("Test #1", func(t *testing.T) {
		cfg := config.New()
		s, err := New(
			srv.URL,
			cfg.ReportInterval,
			metricBuff,
			&key,
			cfg.RateLimit,
			"./testdata/crypto/public.crt",
		)
		require.NoError(t, err)
		err = s.SendQueryUpdateMetrics()
		assert.NoError(t, err)
	})
}
