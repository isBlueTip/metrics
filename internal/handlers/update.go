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

//func UpdateMetric(storage repository.StorageInterface) http.HandlerFunc {
//	return func(res http.ResponseWriter, req *http.Request) {
//		log.Println("")
//		log.Printf("   %s: %v   \n", req.Method, req.URL.Path)
//		if req.Method != http.MethodPost {
//			http.Error(res, "only POST request accepted", http.StatusMethodNotAllowed)
//			return
//		}
//
//		splitURL := strings.Split(req.URL.Path, "/")[1:]
//
//		// validate
//		if len(splitURL) < 3 {
//			http.Error(res, "not enough URL params to write metric, expected 3", http.StatusNotFound)
//			return
//		}
//		parsedValue, err := ParseValue(splitURL[0], splitURL[2])
//		if err != nil {
//			if errors.Is(err, strconv.ErrSyntax) {
//				http.Error(res, "invalid metric value, expected numeric", http.StatusBadRequest)
//			} else {
//				http.Error(res, err.Error(), http.StatusBadRequest)
//			}
//			return
//		}
//
//		switch splitURL[0] {
//		case models.Gauge:
//			service.UpdateGauge(storage, splitURL[1], parsedValue.(float64))
//		case models.Counter:
//			service.UpdateCounter(storage, splitURL[1], parsedValue.(int64))
//		}
//
//		body := fmt.Sprintf("struct: %+v\n", storage)
//
//		res.Write([]byte(body))
//	}
//}

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

		//body := fmt.Sprintf("struct: %+v\n", storage)
		body := fmt.Sprint("{}")

		res.Write([]byte(body))
	}
}
