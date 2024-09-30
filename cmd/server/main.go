package main

import (
	"net/http"

	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers/storage/update"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memStorage"
)

var storage *memStorage.Storage

func main() {

	storage = memStorage.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/{typeMetric}/{nameMetric}/{valueMetric}", update.New(storage))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
