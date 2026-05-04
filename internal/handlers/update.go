package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/isBlueTip/metrics/internal/logger"
	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/service"
	"go.uber.org/zap"
)

func UpdateURL(storage repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := strings.ToLower(chi.URLParam(req, "metricType"))
		metricName := strings.ToLower(chi.URLParam(req, "metricName"))
		metricVal := chi.URLParam(req, "metricVal")

		switch metricType {
		case models.Gauge:
			parsedValue, err := strconv.ParseFloat(metricVal, 64)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			logger.Log.Debug(fmt.Sprintf("storage type: %+v\n", storage))
			err = service.UpdateGauge(storage, metricName, parsedValue)
			if err != nil {
				logger.Log.Error("updating error", zap.String("handler", err.Error()))
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			logger.Log.Debug("gauge updated")
		case models.Counter:
			parsedValue, err := strconv.ParseInt(metricVal, 10, 64)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			err = service.UpdateCounter(storage, metricName, parsedValue)
			if err != nil {
				logger.Log.Error("updating error", zap.String("handler", err.Error()))
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(res, fmt.Sprintf(
				"unknown metric type: %s, expected '%s' or '%s'",
				metricType, models.Gauge, models.Counter), http.StatusBadRequest,
			)
		}

		body := "{}"

		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		res.Write([]byte(body))
	}
}

func UpdateJSON(storage repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		body, err := getReader(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer body.Close()

		decoder := json.NewDecoder(body)
		var metric models.Update
		err = decoder.Decode(&metric)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case models.Gauge:
			if metric.Value == nil {
				err = fmt.Errorf("no value provided")
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			metric.Delta = nil
			err = service.UpdateGauge(storage, metric.ID, *metric.Value)
			if err != nil {
				logger.Log.Error("updating error", zap.String("handler", err.Error()))
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		case models.Counter:
			if metric.Delta == nil {
				err = fmt.Errorf("no delta provided")
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			metric.Value = nil
			err = service.UpdateCounter(storage, metric.ID, *metric.Delta)
			if err != nil {
				logger.Log.Error("updating error", zap.String("handler", err.Error()))
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			err = fmt.Errorf("unknown metric type: %s, expected '%s' or '%s'",
				metric.MType, models.Gauge, models.Counter)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(res)
		err = encoder.Encode(&metric)
		if err != nil {
			logger.Log.Error("server error", zap.String("update", err.Error()))
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdateBatch(storage repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		body, err := getReader(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer body.Close()

		decoder := json.NewDecoder(body)

		var metrics []models.Update
		err = decoder.Decode(&metrics)
		if err != nil {
			logger.Log.Error("decoding error", zap.String("handler", err.Error()))
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		validMetrics := make([]models.Update, 0, len(metrics))
		for _, metric := range metrics {
			switch metric.MType {
			case models.Gauge:
				if metric.Value == nil {
					http.Error(res, "no value provided for gauge", http.StatusBadRequest)
					return
				}
				validMetrics = append(validMetrics, metric)
			case models.Counter:
				if metric.Delta == nil {
					http.Error(res, "no delta provided for counter", http.StatusBadRequest)
					return
				}
				validMetrics = append(validMetrics, metric)
			default:
				http.Error(res, "unknown metric type", http.StatusBadRequest)
				return
			}
		}

		err = service.UpdateBatch(storage, validMetrics)
		if err != nil {
			logger.Log.Error("updating error", zap.String("handler", err.Error()))
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(res)
		err = encoder.Encode(metrics)
		if err != nil {
			logger.Log.Error("encoding error", zap.String("handler", err.Error()))
			return
		}
	}
}
