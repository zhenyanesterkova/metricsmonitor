package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handler"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
	"github.com/zhenyanesterkova/metricsmonitor/web"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, reqBody string) (*http.Response, string) {
	t.Helper()

	var buff bytes.Buffer
	_, err := buff.WriteString(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest(method, ts.URL+path, &buff)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	return resp, string(respBody)
}

func CreateTestMemStorage() (storage *memstorage.MemStorage) {
	storage = memstorage.New()

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

func getExpectedHTML(templateName, nameData string, data interface{}) (string, error) {
	buf := bytes.NewBufferString("")

	tmpl, err := template.ParseFS(web.Templates, "template/"+templateName)
	if err != nil {
		return "", fmt.Errorf("main_test.go func getExpectedHTML(): error parse FS - %w", err)
	}
	err = tmpl.ExecuteTemplate(buf, nameData, data)
	if err != nil {
		return "", fmt.Errorf("main_test.go func getExpectedHTML(): error in Execute Template - %w", err)
	}
	return buf.String(), nil
}

func TestRouter(t *testing.T) {
	respHTML, _ := getExpectedHTML("allMetricsView.html", "metrics", [][2]string{
		{"testCounter", "1"},
		{"testGauge", "2.5"},
	})

	memStorage := CreateTestMemStorage()

	loggerInst := logger.NewLogrusLogger()
	err := loggerInst.SetLevelForLog("debug")
	require.NoError(t, err)

	router := chi.NewRouter()

	repoHandler := handler.NewRepositorieHandler(memStorage, loggerInst, nil)
	repoHandler.InitChiRouter(router)

	ts := httptest.NewServer(router)

	defer ts.Close()

	tests := []struct {
		name                          string
		method                        string
		url                           string
		reqBody                       string
		wantRespBody                  string
		metricName                    string
		wantStorageCounterMetricValue int64
		wantStorageGaugeMetricValue   float64
		status                        int
	}{
		{
			name:         "test #1: GET / ",
			method:       http.MethodGet,
			url:          "/",
			wantRespBody: respHTML,
			status:       http.StatusOK,
		},
		{
			name:         "test #2: GET /value/counter/test",
			method:       http.MethodGet,
			url:          "/value/counter/test",
			wantRespBody: "",
			status:       http.StatusNotFound,
		},
		{
			name:         "test #3: GET /value/counter/testCounter",
			method:       http.MethodGet,
			url:          "/value/counter/testCounter",
			wantRespBody: "1",
			status:       http.StatusOK,
		},
		{
			name:         "test #4: GET /value/gauge/testCounter",
			method:       http.MethodGet,
			url:          "/value/gauge/testCounter",
			wantRespBody: "",
			status:       http.StatusNotFound,
		},
		{
			name:                          "test #5: POST /update/counter/testCounter/1",
			method:                        http.MethodPost,
			url:                           "/update/counter/testCounter/1",
			wantRespBody:                  "",
			metricName:                    "testCounter",
			wantStorageCounterMetricValue: int64(2),
			status:                        http.StatusOK,
		},
		{
			name:                          "test #6: POST /update/counter/testCounter/ttt",
			method:                        http.MethodPost,
			url:                           "/update/counter/testCounter/ttt",
			wantRespBody:                  "",
			metricName:                    "testCounter",
			wantStorageCounterMetricValue: int64(2),
			status:                        http.StatusBadRequest,
		},
		{
			name:                          "test #7: POST /update/gauge/testCounter/1",
			method:                        http.MethodPost,
			url:                           "/update/gauge/testCounter/1",
			wantRespBody:                  "",
			metricName:                    "testCounter",
			wantStorageCounterMetricValue: int64(2),
			status:                        http.StatusBadRequest,
		},
		{
			name:                          "test #8: POST /update/counter/testCounterNew/1",
			method:                        http.MethodPost,
			url:                           "/update/counter/testCounterNew/1",
			wantRespBody:                  "",
			metricName:                    "testCounterNew",
			wantStorageCounterMetricValue: int64(1),
			status:                        http.StatusOK,
		},
		{
			name:                        "test #9: POST /update/gauge/testGauge/1.5",
			method:                      http.MethodPost,
			url:                         "/update/gauge/testGauge/1.5",
			wantRespBody:                "",
			metricName:                  "testGauge",
			wantStorageGaugeMetricValue: 1.5,
			status:                      http.StatusOK,
		},
		{
			name:                        "test #10: POST /update/gauge/testGauge/ttt",
			method:                      http.MethodPost,
			url:                         "/update/gauge/testGauge/ttt",
			wantRespBody:                "",
			metricName:                  "testGauge",
			wantStorageGaugeMetricValue: 1.5,
			status:                      http.StatusBadRequest,
		},
		{
			name:                        "test #11: POST /update/counter/testGauge/1",
			method:                      http.MethodPost,
			url:                         "/update/counter/testGauge/1",
			wantRespBody:                "",
			metricName:                  "testGauge",
			wantStorageGaugeMetricValue: 1.5,
			status:                      http.StatusBadRequest,
		},
		{
			name:                        "test #12: POST /update/gauge/testGaugeNew/3.6",
			method:                      http.MethodPost,
			url:                         "/update/gauge/testGaugeNew/3.6",
			wantRespBody:                "",
			metricName:                  "testGaugeNew",
			wantStorageGaugeMetricValue: 3.6,
			status:                      http.StatusOK,
		},
		{
			name:         "test #13: GET /value/gauge/test",
			method:       http.MethodGet,
			url:          "/value/gauge/test",
			wantRespBody: "",
			status:       http.StatusNotFound,
		},
		{
			name:         "test #14: GET /value/gauge/testGauge",
			method:       http.MethodGet,
			url:          "/value/gauge/testGauge",
			wantRespBody: "1.5",
			status:       http.StatusOK,
		},
		{
			name:         "test #15: GET /update/gauge/testGaugeNew/3",
			method:       http.MethodGet,
			url:          "/update/gauge/testGaugeNew/3.6",
			wantRespBody: "",
			status:       http.StatusMethodNotAllowed,
		},
		{
			name:         "test #16: POST /value/gauge/testGaugeNew",
			method:       http.MethodPost,
			url:          "/value/gauge/testGaugeNew",
			wantRespBody: "",
			status:       http.StatusMethodNotAllowed,
		},
		{
			name:         "test #17: GET /ping",
			method:       http.MethodGet,
			url:          "/ping",
			wantRespBody: "",
			status:       http.StatusOK,
		},
		{
			name:         "test #18: POST /value/",
			method:       http.MethodPost,
			url:          "/value/",
			reqBody:      `{"type": "gauge", "id": "testGauge"}`,
			wantRespBody: "{\"value\":1.5,\"id\":\"testGauge\",\"type\":\"gauge\"}\n",
			status:       http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, respBody := testRequest(t, ts, test.method, test.url, test.reqBody)

			assert.Equal(t, test.status, resp.StatusCode)
			assert.Equal(t, test.wantRespBody, respBody)

			if test.wantStorageCounterMetricValue != 0 {
				actualValue, err := memStorage.GetMetricValue(test.metricName, "counter")
				require.NoError(t, err)
				assert.Equal(t, test.wantStorageCounterMetricValue, *actualValue.Delta)
			}
			if test.wantStorageGaugeMetricValue != 0.0 {
				actualValue, err := memStorage.GetMetricValue(test.metricName, "gauge")
				require.NoError(t, err)
				assert.Equal(t, test.wantStorageGaugeMetricValue, *actualValue.Value)
			}
			err := resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
