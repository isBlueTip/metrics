package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/isBlueTip/metrics/internal/handlers"
	"github.com/isBlueTip/metrics/internal/hash"
	"github.com/isBlueTip/metrics/internal/logger"
	"github.com/isBlueTip/metrics/internal/repository"
	"go.uber.org/zap"
)

func Router(storage repository.Storage, DBConn string, key string) *chi.Mux {
	r := chi.NewRouter()

	// Recovery middleware to catch panics
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Log.Error("panic recovered", zap.Any("error", err))
					http.Error(w, fmt.Sprintf("internal error: %v", err), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})

	// Hash middleware for request verification and response signing
	if key != "" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedHash := r.Header.Get("HashSHA256")

				// Verify request hash if header is present
				if receivedHash != "" {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					r.Body.Close()

					expectedHash := hash.Compute(body, key)
					if expectedHash != receivedHash {
						logger.Log.Error("hash mismatch",
							zap.String("expected", expectedHash),
							zap.String("received", receivedHash))
						http.Error(w, "hash mismatch", http.StatusBadRequest)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(body))
				}

				// Wrap response writer to capture body for response hash
				wr := &hashResponseWriter{ResponseWriter: w}

				next.ServeHTTP(wr, r)

				// Add hash to response
				if wr.body.Len() > 0 {
					h := hash.Compute(wr.body.Bytes(), key)
					w.Header().Set("HashSHA256", h)
				}
			})
		})
	}

	//r.Use(middleware.Logger)
	//r.Use(logger.Log)
	r.Use(GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(storage))
		r.Post("/update/{metricType}/{metricName:[a-zA-Z0-9_]+}/{metricVal}", logger.RequestLogger(handlers.UpdateURL(storage)))
		r.Post("/update/", logger.RequestLogger(handlers.UpdateJSON(storage)))
		r.Post("/updates/", logger.RequestLogger(handlers.UpdateBatch(storage)))
		r.Get("/value/{metricType}/{metricName:[a-zA-Z0-9_]+}", logger.RequestLogger(handlers.GetByNameURL(storage)))
		r.Post("/value/", logger.RequestLogger(handlers.GetByNameJSON(storage)))
		r.Get("/ping", logger.RequestLogger(handlers.Ping(storage)))
	})

	return r
}

type hashResponseWriter struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (w *hashResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
