package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
	"github.com/zhenyanesterkova/metricsmonitor/web"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)

}

func CreateTestMemStorage() (storage *memstorage.MemStorage) {
	storage = memstorage.New()
	_ = storage.UpdateMetric("testCounter", "counter", "1")
	_ = storage.UpdateMetric("testGauge", "gauge", "2.5")
	return
}

func getExpectedHTML(templateName, nameData string, data interface{}) (string, error) {
	buf := bytes.NewBufferString("")

	tmpl, err := template.ParseFS(web.Templates, "template/"+templateName)
	if err != nil {
		return "", err
	}
	err = tmpl.ExecuteTemplate(buf, nameData, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestRouter(t *testing.T) {

	respHTML, _ := getExpectedHTML("allMetricsView.html", "metrics", [][2]string{
		{"testCounter", "1"},
		{"testGauge", "2.5"},
	})

	memStorage := CreateTestMemStorage()

	router := chi.NewRouter()

	repoHandler := handlers.NewRepositorieHandler(memStorage)
	repoHandler.InitChiRouter(router)

	ts := httptest.NewServer(router)

	defer ts.Close()

	tests := []struct {
		name                          string
		method                        string
		url                           string
		wantRespBody                  string
		metricName                    string
		wantStorageCounterMetricValue string
		wantStorageGaugeMetricValue   string
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
			name:         "test #2: method GET /value/counter/test (unknown metric name)",
			method:       http.MethodGet,
			url:          "/value/counter/test",
			wantRespBody: "",
			status:       http.StatusNotFound,
		},
		{
			name:         "test #3: method GET /value/counter/testCounter (correct metric name, type)",
			method:       http.MethodGet,
			url:          "/value/counter/testCounter",
			wantRespBody: "1",
			status:       http.StatusOK,
		},
		{
			name:         "test #4: method GET /value/gauge/testCounter (correct metric name, incorrect metric type)",
			method:       http.MethodGet,
			url:          "/value/gauge/testCounter",
			wantRespBody: "",
			status:       http.StatusNotFound,
		},
		{
			name:                          "test #5: method POST /update/counter/testCounter/1 (correct metric name, value, type; existing metric; counter)",
			method:                        http.MethodPost,
			url:                           "/update/counter/testCounter/1",
			wantRespBody:                  "",
			metricName:                    "testCounter",
			wantStorageCounterMetricValue: "2",
			status:                        http.StatusOK,
		},
		{
			name:                          "test #6: method POST /update/counter/testCounter/ttt (correct metric name, type; incorrect value; existing metric; counter)",
			method:                        http.MethodPost,
			url:                           "/update/counter/testCounter/ttt",
			wantRespBody:                  "",
			metricName:                    "testCounter",
			wantStorageCounterMetricValue: "2",
			status:                        http.StatusBadRequest,
		},
		{
			name:                          "test #7: method POST /update/gauge/testCounter/1 (correct metric name, value; incorrect type; existing metric; counter)",
			method:                        http.MethodPost,
			url:                           "/update/gauge/testCounter/1",
			wantRespBody:                  "",
			metricName:                    "testCounter",
			wantStorageCounterMetricValue: "2",
			status:                        http.StatusBadRequest,
		},
		{
			name:                          "test #8: method POST /update/counter/testCounterNew/1 (correct metric name, value, type; not existing metric; counter)",
			method:                        http.MethodPost,
			url:                           "/update/counter/testCounterNew/1",
			wantRespBody:                  "",
			metricName:                    "testCounterNew",
			wantStorageCounterMetricValue: "1",
			status:                        http.StatusOK,
		},
		{
			name:                        "test #9: method POST /update/gauge/testGauge/1.5 (correct metric name, value, type; existing metric; gauge)",
			method:                      http.MethodPost,
			url:                         "/update/gauge/testGauge/1.5",
			wantRespBody:                "",
			metricName:                  "testGauge",
			wantStorageGaugeMetricValue: "1.5",
			status:                      http.StatusOK,
		},
		{
			name:                        "test #10: method POST /update/gauge/testGauge/ttt (correct metric name, type; incorrect value; existing metric; gauge)",
			method:                      http.MethodPost,
			url:                         "/update/gauge/testGauge/ttt",
			wantRespBody:                "",
			metricName:                  "testGauge",
			wantStorageGaugeMetricValue: "1.5",
			status:                      http.StatusBadRequest,
		},
		{
			name:                        "test #11: method POST /update/counter/testGauge/1 (correct metric name, value; incorrect type; existing metric; gauge)",
			method:                      http.MethodPost,
			url:                         "/update/counter/testGauge/1",
			wantRespBody:                "",
			metricName:                  "testGauge",
			wantStorageGaugeMetricValue: "1.5",
			status:                      http.StatusBadRequest,
		},
		{
			name:                        "test #12: method POST /update/gauge/testGaugeNew/3.6 (correct metric name, value, type; not existing metric; gauge)",
			method:                      http.MethodPost,
			url:                         "/update/gauge/testGaugeNew/3.6",
			wantRespBody:                "",
			metricName:                  "testGaugeNew",
			wantStorageGaugeMetricValue: "3.6",
			status:                      http.StatusOK,
		},
		{
			name:         "test #13: method GET /value/gauge/test (unknown metric name)",
			method:       http.MethodGet,
			url:          "/value/gauge/test",
			wantRespBody: "",
			status:       http.StatusNotFound,
		},
		{
			name:         "test #14: method GET /value/gauge/testGauge (correct metric name, type)",
			method:       http.MethodGet,
			url:          "/value/gauge/testGauge",
			wantRespBody: "1.5",
			status:       http.StatusOK,
		},
		{
			name:         "test #15: method GET /update/gauge/testGaugeNew/3 (correct metric name, type, value; incorrect method)",
			method:       http.MethodGet,
			url:          "/update/gauge/testGaugeNew/3.6",
			wantRespBody: "",
			status:       http.StatusMethodNotAllowed,
		},
		{
			name:         "test #16: method POST /value/gauge/testGaugeNew (correct metric name, type, value; incorrect method)",
			method:       http.MethodPost,
			url:          "/value/gauge/testGaugeNew",
			wantRespBody: "",
			status:       http.StatusMethodNotAllowed,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, respBody := testRequest(t, ts, test.method, test.url)
			defer resp.Body.Close()

			assert.Equal(t, test.status, resp.StatusCode)
			assert.Equal(t, test.wantRespBody, respBody)

			if test.wantStorageCounterMetricValue != "" {
				actualValue, err := memStorage.GetMetricValue(test.metricName, "counter")
				require.NoError(t, err)
				assert.Equal(t, test.wantStorageCounterMetricValue, actualValue)
			}
			if test.wantStorageGaugeMetricValue != "" {
				actualValue, err := memStorage.GetMetricValue(test.metricName, "gauge")
				require.NoError(t, err)
				assert.Equal(t, test.wantStorageGaugeMetricValue, actualValue)
			}
		})
	}
}
