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

		decoder := json.NewDecoder(body)
		defer req.Body.Close()

		var metrics models.Metrics
		err = decoder.Decode(&metrics)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		switch metrics.MType {
		case models.Gauge:
			if metrics.Value == nil {
				err = fmt.Errorf("no value provided")
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			metrics.Delta = nil
			service.UpdateGauge(storage, metrics.ID, *metrics.Value)
		case models.Counter:
			if metrics.Delta == nil {
				err = fmt.Errorf("no delta provided")
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			metrics.Value = nil
			service.UpdateCounter(storage, metrics.ID, *metrics.Delta)
		default:
			err = fmt.Errorf("unknown metric type: %s, expected '%s' or '%s'",
				metrics.MType, models.Gauge, models.Counter)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		encoder := json.NewEncoder(res)
		err = encoder.Encode(&metrics)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
