package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/service"
)

type URLMetric struct {
	Type, Name string
	Value      interface{}
}

func UpdateMetric(storage repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		metricVal := chi.URLParam(req, "metricVal")

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

		res.Write([]byte(body))
	}
}
