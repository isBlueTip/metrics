package server

import (
	"github.com/isbluetip/metrics/internal/repository"
	"github.com/xxx/metrics/internal/handlers"
	"net/http"
)

func Router(storage repository.StorageInterface) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update", handlers.UpdateMetric(storage)))
	return mux
}
