package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/service"
)

func UpdateURL(storage repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := strings.ToLower(chi.URLParam(req, "metricType"))
		metricName := strings.ToLower(chi.URLParam(req, "metricName"))
		metricVal := chi.URLParam(req, "metricVal")

		// todo branch depending on metricType here and remove interface{}
		parsedValue, err := ParseValue(metricType, metricVal)
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				http.Error(res, "invalid metric value, expected numeric", http.StatusBadRequest)
			} else {
				http.Error(res, err.Error(), http.StatusBadRequest)
			}
			return
		}

		switch metricType {
		case models.Gauge:
			service.UpdateGauge(storage, metricName, parsedValue.(float64))
		case models.Counter:
			service.UpdateCounter(storage, metricName, parsedValue.(int64))
		}

		body := "{}"

		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		res.Write([]byte(body))
	}
}

func UpdateJSON(storage repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(req.Body)
		defer req.Body.Close()

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var metrics Metrics
		err = json.Unmarshal(buf.Bytes(), &metrics)

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
			err = fmt.Errorf("unknown metric type: %s, expected '%s' or '%s'", metrics.MType, models.Gauge, models.Counter)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(metrics)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(body)
	}
}
