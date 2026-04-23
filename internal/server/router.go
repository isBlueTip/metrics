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
	r.Use(GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(storage))
		r.Post("/update/{metricType}/{metricName:[a-zA-Z0-9_]+}/{metricVal}", logger.RequestLogger(handlers.UpdateURL(storage)))
		r.Post("/update/", logger.RequestLogger(handlers.UpdateJSON(storage)))
		r.Get("/value/{metricType}/{metricName:[a-zA-Z0-9_]+}", logger.RequestLogger(handlers.GetByNameURL(storage)))
		r.Post("/value/", logger.RequestLogger(handlers.GetByNameJSON(storage)))
	})

	return r
}
