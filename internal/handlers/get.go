package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/service"
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
		gauges := service.GetGauges(s)
		counters := service.GetCounters(s)

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
			panic(err)
		}

		data := Data{
			Title:   "title",
			Metrics: template.HTML(strings.Join(buf, "\n")),
		}

		err = tmpl.Execute(res, data)
		if err != nil {
			panic(err)
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
			val, err := service.GetGauge(s, metrics.ID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			metrics.Delta = nil
			metrics.Value = &val
		case models.Counter:
			val, err := service.GetCounter(s, metrics.ID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			metrics.Value = nil
			metrics.Delta = &val
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
