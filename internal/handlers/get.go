package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/service"
)

func GetAllMetrics(s repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var buf []string
		gauges := service.GetGauges(s)
		counters := service.GetCounters(s)

		for _, g := range gauges {
			str := fmt.Sprintf("%s: %.2f\n", g.Name, g.Value)
			buf = append(buf, str)
		}

		for _, c := range counters {
			str := fmt.Sprintf("%s: %d\n", c.Name, c.Value)
			buf = append(buf, str)
		}
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte(strings.Join(buf, "")))
	}
}

func GetMetricByName(s repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		metricType := strings.ToLower(chi.URLParam(req, "metricType"))
		metricName := strings.ToLower(chi.URLParam(req, "metricName"))
		var buf []byte

		switch metricType {

		case models.Gauge:
			val, err := service.GetGauge(s, metricName)
			if err != nil {
				http.Error(res, "", http.StatusNotFound)
				return
			}
			buf = strconv.AppendFloat(buf, val, 'f', -1, 64)

		case models.Counter:
			val, err := service.GetCounter(s, metricName)
			if err != nil {
				http.Error(res, "", http.StatusNotFound)
				return
			}
			buf = strconv.AppendInt(buf, val, 10)

		default:
			http.Error(res, fmt.Sprintf("unknown metric type: %s, expected '%s' or '%s'", metricType, models.Gauge, models.Counter), http.StatusBadRequest)
			return
		}
		res.Write(buf)
	}
}
