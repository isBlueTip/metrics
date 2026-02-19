package handlers

import (
	"errors"
	"fmt"
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

func UpdateMetric(storage repository.StorageInterface) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := chi.URLParam(req, "metricType")
		metricName := chi.URLParam(req, "metricName")
		metricVal := chi.URLParam(req, "metricVal")

		if metricType != models.Gauge && metricType != models.Counter {
			err := fmt.Errorf("metricType: '%s', expected '%s' or '%s'", metricType, models.Gauge, models.Counter)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

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
