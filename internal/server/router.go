package server

import (
	"github.com/xxx/metrics/internal/handlers"
	"github.com/xxx/metrics/internal/repository"
	"net/http"
)

func Router(storage repository.StorageInterface) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update", handlers.UpdateMetric(storage)))
	return mux
}
