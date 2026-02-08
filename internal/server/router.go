package server

import (
	//"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/isBlueTip/metrics/internal/handlers"
	"github.com/isBlueTip/metrics/internal/repository"
)

// func Router(storage repository.StorageInterface) *http.ServeMux {
func Router(storage repository.StorageInterface) *chi.Mux {
	//mux := http.NewServeMux()
	//mux.Handle("/update/", http.StripPrefix("/update", handlers.UpdateMetric(storage)))
	//return mux

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Post("/update/*", handlers.UpdateMetric(storage))
	r.Post("/update/{metricType}/{metricName}/{metricVal}", handlers.UpdateMetric(storage))

	return r
}
