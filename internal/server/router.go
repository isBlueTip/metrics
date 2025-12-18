package server

import (
	"github.com/xxx/metrics/internal/handler"
	"github.com/xxx/metrics/internal/repository"
	"net/http"
)

func Router(storage repository.StorageInterface) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update", handler.UpdateMetric(storage)))
	return mux
}
