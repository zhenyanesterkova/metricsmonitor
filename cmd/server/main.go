package main

import (
	"net/http"

	"github.com/zhenyanesterkova/metricsmonitor/internal/handlers/storage/update"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

var storage *memstorage.Storage

func main() {

	storage = memstorage.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/{typeMetric}/{nameMetric}/{valueMetric}", update.New(storage))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
