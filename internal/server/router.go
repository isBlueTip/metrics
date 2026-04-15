package server

import (
	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/isBlueTip/metrics/internal/handlers"
	//"github.com/isBlueTip/metrics/internal/logger"
	"github.com/isBlueTip/metrics/internal/logger"
	"github.com/isBlueTip/metrics/internal/repository"
	//"go.uber.org/zap"
)

func Router(storage repository.Storage) *chi.Mux {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	//r.Use(logger.Log)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAllMetrics(storage))
		r.Post("/update/{metricType}/{metricName:[a-zA-Z0-9_]+}/{metricVal}", logger.RequestLogger(handlers.UpdateMetric(storage)))
		r.Get("/value/{metricType}/{metricName:[a-zA-Z0-9_]+}", logger.RequestLogger(handlers.GetMetricByName(storage)))
	})

	return r
}
