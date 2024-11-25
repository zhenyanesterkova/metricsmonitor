package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/metric"
)

func (rh *RepositorieHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "typeMetric")
	metricName := chi.URLParam(r, "nameMetric")
	metricValue := chi.URLParam(r, "valueMetric")

	metrica := metric.New(metricType)
	metrica.ID = metricName
	switch metricType {
	case metric.TypeGauge:
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		*metrica.Value = val
	case metric.TypeCounter:
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		*metrica.Delta = val
	}

	_, err := rh.Repo.UpdateMetric(metrica)
	if err != nil {
		switch {
		case errors.Is(err, metric.ErrInvalidName):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, metric.ErrParseValue) ||
			errors.Is(err, metric.ErrUnknownType) ||
			errors.Is(err, metric.ErrInvalidType):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rh *RepositorieHandler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	log.Info("updating metric ...")

	newMetric := metric.New("")
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&newMetric); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	log.WithFields(logrus.Fields{
		"ID":    newMetric.ID,
		"Type":  newMetric.MType,
		"Value": *newMetric.Value,
		"Delta": *newMetric.Delta,
	}).Info("metric for updating")

	updating, err := rh.Repo.UpdateMetric(newMetric)
	if err != nil {
		switch {
		case errors.Is(err, metric.ErrInvalidName):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, metric.ErrParseValue) ||
			errors.Is(err, metric.ErrUnknownType) ||
			errors.Is(err, metric.ErrInvalidType):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(updating); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}
