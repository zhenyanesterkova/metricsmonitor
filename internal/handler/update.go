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
	log := rh.Logger.LogrusLog

	log.Info("updating metric ...")

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
		log.WithFields(logrus.Fields{
			"ID":    metrica.ID,
			"Type":  metrica.MType,
			"Value": *metrica.Value,
		}).Info("gauge metric updating")
	case metric.TypeCounter:
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		*metrica.Delta = val
		log.WithFields(logrus.Fields{
			"ID":    metrica.ID,
			"Type":  metrica.MType,
			"Delta": *metrica.Delta,
		}).Info("counter metric updating")
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
			log.Errorf("handler func UpdateMetric(): error update metric - %v", err)
			http.Error(w, TextServerError, http.StatusInternalServerError)
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
		log.Errorf("handler func UpdateMetricJSON(): error decode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
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
			log.Errorf("handler func UpdateMetricJSON(): error update metric - %v", err)
			http.Error(w, TextServerError, http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(updating); err != nil {
		log.Errorf("handler func UpdateMetricJSON(): error encode metric - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}
}

func (rh *RepositorieHandler) UpdateManyMetrics(w http.ResponseWriter, r *http.Request) {
	log := rh.Logger.LogrusLog

	metricsList := []metric.Metric{}
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&metricsList); err != nil {
		log.Errorf("handler func UpdateManyMetrics(): error decode metrics - %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	err := rh.Repo.UpdateManyMetrics(r.Context(), metricsList)
	if err != nil {
		log.Errorf("handler func UpdateManyMetrics(): error update metrics: %v", err)
		http.Error(w, TextServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
