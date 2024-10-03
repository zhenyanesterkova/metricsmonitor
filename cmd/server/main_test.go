package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers/storage/update"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

func TestUpdateHandler(t *testing.T) {
	type path struct {
		typeMetric  string
		nameMetric  string
		valueMetric string
	}
	type want struct {
		statuseCode int
		contentType string
	}
	tests := []struct {
		name    string
		method  string
		request string
		path    path
		want    want
	}{
		{
			name:    "test #1: method GET",
			method:  http.MethodGet,
			request: "/update/counter/test/1",
			path: path{
				typeMetric:  "counter",
				nameMetric:  "test",
				valueMetric: "1",
			},
			want: want{
				statuseCode: http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
		{
			name:    "test #2: method POST without name",
			method:  http.MethodPost,
			request: "/update/counter//1",
			path: path{
				typeMetric:  "counter",
				nameMetric:  "",
				valueMetric: "1",
			},
			want: want{
				statuseCode: http.StatusNotFound,
				contentType: "",
			},
		},
		{
			name:    "test #3: method POST incorrect type",
			method:  http.MethodPost,
			request: "/update/ttt/test/1",
			path: path{
				typeMetric:  "ttt",
				nameMetric:  "test",
				valueMetric: "1",
			},
			want: want{
				statuseCode: http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name:    "test #4: method POST incorrect value",
			method:  http.MethodPost,
			request: "/update/counter/test/",
			path: path{
				typeMetric:  "counter",
				nameMetric:  "test",
				valueMetric: "",
			},
			want: want{
				statuseCode: http.StatusBadRequest,
				contentType: "",
			},
		},
		{
			name:    "test #4: method POST correct",
			method:  http.MethodPost,
			request: "/update/counter/test/1",
			path: path{
				typeMetric:  "counter",
				nameMetric:  "test",
				valueMetric: "1",
			},
			want: want{
				statuseCode: http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newStorage := memstorage.New()
			req := httptest.NewRequest(test.method, test.request, nil)
			req.SetPathValue("typeMetric", test.path.typeMetric)
			req.SetPathValue("nameMetric", test.path.nameMetric)
			req.SetPathValue("valueMetric", test.path.valueMetric)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(update.New(newStorage))
			h(w, req)

			result := w.Result()

			assert.Equal(t, test.want.statuseCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
