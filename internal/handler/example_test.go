package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func ExampleRepositorieHandler_UpdateMetricJSON() {
	type Metric struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Value *float64 `json:"value,omitempty"`
	}
	// Create a new HTTP client
	client := &http.Client{}

	// Prepare the metric data to be updated
	metricData := Metric{
		ID:    "Alloc",
		MType: "gauge",
		Value: new(float64),
	}
	*metricData.Value = 23.4

	// Marshal the metric data into JSON
	jsonData, err := json.Marshal(metricData)
	if err != nil {
		panic(err)
	}

	// Create a new request
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/update/",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
}

func ExampleRepositorieHandler_UpdateMetric() {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/update/counter/PollCount/3",
		http.NoBody,
	)
	if err != nil {
		panic(err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
}

func ExampleRepositorieHandler_UpdateManyMetrics() {
	type Metric struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Value *float64 `json:"value,omitempty"`
		Delta *int64   `json:"delta,omitempty"`
	}
	// Create a new HTTP client
	client := &http.Client{}

	// Prepare the metric data to be updated
	metricDataGauge := Metric{
		ID:    "Alloc",
		MType: "gauge",
		Value: new(float64),
	}
	*metricDataGauge.Value = 23.4
	metricDataCounter := Metric{
		ID:    "PollCount",
		MType: "counter",
		Delta: new(int64),
	}
	*metricDataCounter.Delta = 23

	reqBody := []Metric{metricDataGauge, metricDataCounter}

	// Marshal the metric data into JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	// Create a new request
	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/update/counter/PollCount/3",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
}
