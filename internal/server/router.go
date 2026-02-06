package server

import (
	"github.com/isBlueTip/metrics/internal/handlers"
	"github.com/isBlueTip/metrics/internal/repository"
	"net/http"
)

func Router(storage repository.StorageInterface) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update", handlers.UpdateMetric(storage)))
	return mux
}
