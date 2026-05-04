package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
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

const htmlTmpl = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<p>{{.Metrics}}</p>
	</body>
	</html>
	`

type Data struct {
	Title   string
	Metrics template.HTML
}

func GetAll(s repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var buf []string
		gauges, err := service.GetGauges(s)
		if err != nil {
			logger.Log.Error("gauges retrieving error", zap.String("gauges", err.Error()))
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		counters, err := service.GetCounters(s)
		if err != nil {
			logger.Log.Error("counters retrieving error", zap.String("counters", err.Error()))
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, g := range gauges {
			str := fmt.Sprintf("%s: %.2f<br>", g.Name, g.Value)
			buf = append(buf, str)
		}

		for _, c := range counters {
			str := fmt.Sprintf("%s: %d<br>", c.Name, c.Value)
			buf = append(buf, str)
		}
		res.Header().Set("Content-Type", "text/html")
		tmpl, err := template.New("webpage").Parse(htmlTmpl)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		data := Data{
			Title:   "title",
			Metrics: template.HTML(strings.Join(buf, "\n")),
		}

		err = tmpl.Execute(res, data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetByNameURL(s repository.Storage) http.HandlerFunc {
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

func GetByNameJSON(s repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		body, err := getReader(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer body.Close()

		decoder := json.NewDecoder(body)

		var metrics models.Update
		err = decoder.Decode(&metrics)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		switch metrics.MType {
		case models.Gauge:
			val, err := service.GetGauge(s, metrics.ID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
			metrics.Delta = nil
			metrics.Value = &val
		case models.Counter:
			val, err := service.GetCounter(s, metrics.ID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
			metrics.Value = nil
			metrics.Delta = &val
		default:
			err = fmt.Errorf("unknown metric type: %s, expected '%s' or '%s'", metrics.MType, models.Gauge, models.Counter)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(res)
		err = encoder.Encode(metrics)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func Ping(s repository.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := service.Ping(s); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		res.WriteHeader(http.StatusOK)
	}
}
