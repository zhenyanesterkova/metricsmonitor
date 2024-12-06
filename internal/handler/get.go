package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/backoff"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
	"github.com/zhenyanesterkova/metricsmonitor/web"
)

func (rh *RepositorieHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	res, err := rh.Repo.GetAllMetrics()
	if err != nil {
		if rh.checkRetry(err) {
			err = rh.retry(func() error {
				res, err = rh.Repo.GetAllMetrics()
				return fmt.Errorf("failed get all metrics: %w", err)
			})
		}
		if err != nil {
			log.Errorf("handler func GetAllMetrics(): error get metrics - %v", err)
			http.Error(w, TextServerError, http.StatusInternalServerError)
			return
		}
	}

	tmplMetrics, err := template.ParseFS(web.Templates, "template/allMetricsView.html")
	if err != nil {
		log.Errorf("handler func GetAllMetrics(): error parse html template - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err = tmplMetrics.ExecuteTemplate(w, "metrics", res)
	if err != nil {
		log.Errorf("handler func GetAllMetrics(): error execute html template - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
}

func (rh *RepositorieHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	name := chi.URLParam(r, "nameMetric")
	metricType := chi.URLParam(r, "typeMetric")

	res, err := rh.Repo.GetMetricValue(name, metricType)
	if err != nil {
		if errors.Is(err, metric.ErrUnknownMetric) || errors.Is(err, metric.ErrInvalidType) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Errorf("handler func GetMetricValue(): error get metric value - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res.String()))
}

func (rh *RepositorieHandler) GetMetricValueJSON(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	metrica := metric.New("")
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrica); err != nil {
		log.Errorf("handler func GetMetricValueJSON(): error decode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	res, err := rh.Repo.GetMetricValue(metrica.ID, metrica.MType)
	if err != nil {
		if errors.Is(err, metric.ErrUnknownMetric) || errors.Is(err, metric.ErrInvalidType) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Errorf("handler func GetMetricValueJSON(): error get metric value - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		log.Errorf("handler func GetMetricValueJSON(): error encode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
}

func (rh *RepositorieHandler) Ping(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog
	ok, err := rh.Repo.Ping()
	if err != nil || !ok {
		log.Errorf("failed ping storage: %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rh *RepositorieHandler) checkRetry(err error) bool {
	var pgErr *pgconn.PgError
	var pgErrConn *pgconn.ConnectError
	res := false
	if errors.As(err, &pgErr) {
		res = pgerrcode.IsConnectionException(pgErr.Code)
	} else if errors.As(err, &pgErrConn) {
		res = true
	}
	return res
}

func (rh *RepositorieHandler) retry(work func() error) error {
	defer rh.backoff.Reset()
	for {
		err := work()

		if err == nil {
			return nil
		}
		if rh.checkRetry(err) {
			var delay time.Duration
			if delay = rh.backoff.Next(); delay == backoff.Stop {
				return err
			}
			time.Sleep(delay)
		}
	}
}
