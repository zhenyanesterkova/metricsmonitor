package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
)

const (
	hashKey = "testHashKey"
)

type headerParams struct {
	acceptEncoding  []string
	accept          string
	contentEncoding string
}

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
	headersParams headerParams,
	reqBody string,
) *http.Response {
	t.Helper()

	var buff bytes.Buffer
	if headersParams.contentEncoding == "gzip" {
		gzWriter := gzip.NewWriter(&buff)
		defer func() {
			err := gzWriter.Close()
			require.NoError(t, err)
		}()
		_, _ = gzWriter.Write([]byte(reqBody))
	} else {
		_, err := buff.WriteString(reqBody)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, ts.URL+path, &buff)
	require.NoError(t, err)

	for _, enc := range headersParams.acceptEncoding {
		req.Header.Add("Accept-Encoding", enc)
	}
	req.Header.Add("Accept", headersParams.accept)
	req.Header.Add("Content-Encoding", headersParams.contentEncoding)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func TestMiddleware(t *testing.T) {
	loggerInst := logger.NewLogrusLogger()
	err := loggerInst.SetLevelForLog("debug")
	require.NoError(t, err)

	key := hashKey
	mdlWare := NewMiddlewareStruct(loggerInst, &key)

	router := chi.NewRouter()

	router.Use(mdlWare.ResetRespDataStruct)
	router.Use(mdlWare.RequestLogger)
	router.Use(mdlWare.CheckSignData)
	router.Use(mdlWare.GZipMiddleware)
	router.Route("/", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success ping req"))
		})
		r.Get("/gzip", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			_, _ = w.Write([]byte("{}"))
		})
		r.Get("/pingempty", func(w http.ResponseWriter, r *http.Request) {
		})
		r.Get("/debug/pprof/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	ts := httptest.NewServer(router)

	defer ts.Close()

	t.Run("ResetRespDataStruct", func(t *testing.T) {
		params := headerParams{
			acceptEncoding:  []string{},
			accept:          "",
			contentEncoding: "",
		}
		resp := testRequest(
			t,
			ts,
			http.MethodGet,
			"/ping",
			params,
			"")
		defer func() {
			err := resp.Body.Close()
			require.NoError(t, err)
		}()
		resp = testRequest(
			t,
			ts,
			http.MethodGet,
			"/pingempty",
			params,
			"")
		require.Equal(t, mdlWare.respData.responseData.size, 0)
		require.Equal(t, mdlWare.respData.responseData.status, 0)

		resp = testRequest(
			t,
			ts,
			http.MethodGet,
			"/ping",
			params,
			"")
		resp = testRequest(t, ts, http.MethodGet, "/debug/pprof/heap", params, "")
		require.Equal(t, mdlWare.respData.responseData.size, 16)
		require.Equal(t, mdlWare.respData.responseData.status, 200)
	})

	t.Run("GZipMiddleware", func(t *testing.T) {
		params := headerParams{
			acceptEncoding:  []string{"gzip"},
			accept:          "application/json",
			contentEncoding: "gzip",
		}
		resp := testRequest(
			t,
			ts,
			http.MethodPost,
			"/gzip",
			params,
			"{\"test\":\"test\"}",
		)

		got := resp.Header.Get("Content-Encoding")
		require.Equal(t, "gzip", got)

		err := resp.Body.Close()
		require.NoError(t, err)
	})
}

func Test_isPprofPath(t *testing.T) {
	pathDebug := "/debug/pprof/heap"
	pathNoDebug := "/ping"
	pathNoDebug1 := "/debug"

	t.Run("is pprof path", func(t *testing.T) {
		answ := isPprofPath(pathDebug)
		require.True(t, answ)
	})
	t.Run("not is pprof path", func(t *testing.T) {
		answ := isPprofPath(pathNoDebug)
		require.False(t, answ)

		answ = isPprofPath(pathNoDebug1)
		require.False(t, answ)
	})
}

func Test_isCompression(t *testing.T) {
	typeCompression := "application/json"
	typeCompression1 := "text/html"
	typeNoCompression := "text/plain"

	t.Run("is compression", func(t *testing.T) {
		answ := isCompression(typeCompression)
		require.True(t, answ)
		answ = isCompression(typeCompression1)
		require.True(t, answ)
	})
	t.Run("not is compression", func(t *testing.T) {
		answ := isCompression(typeNoCompression)
		require.False(t, answ)
	})
}
