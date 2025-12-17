package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type URLMetric struct {
	Type, Name, Value string
}

func serveMetrics(storage *MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Println("   **********************   ")
		log.Println("")
		if req.Method != http.MethodPost {
			http.Error(res, "only POST request accepted", http.StatusMethodNotAllowed)
			return
		}

		if err := req.ParseForm(); err != nil {
			res.Write([]byte(err.Error()))
			return
		}

		log.Printf("'%v', type: %T\r\n", req.URL.Path, req.URL.Path)
		splitURL := strings.Split(req.URL.Path, "/")[1:]
		//if err := ValidateURL(splitURL); err != nil {
		//	http.Error(res, err.Error(), http.StatusBadRequest)
		//	return
		//}

		if len(splitURL) < 4 {
			http.Error(res, "not enough URL params to write metric, expected 4", http.StatusNotFound)
			return
		}
		if splitURL[0] != "update" {
			http.Error(res, "unknown verb from URL, expected 'update'", http.StatusBadRequest)
			return
		}
		if _, err := strconv.ParseFloat(splitURL[3], 64); err != nil {
			http.Error(res, "invalid metric value, expected numeric", http.StatusBadRequest)
			return
		}

		metric := URLMetric{Type: splitURL[1], Name: splitURL[2], Value: splitURL[3]}

		switch metric.Type {
		case "gauge":
			val, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				http.Error(res, "unable to parse value as float64", http.StatusBadRequest)
				return
			}
			storage.AddGauge(metric.Name, val)
		case "counter":
			val, err := strconv.Atoi(metric.Value)
			if err != nil {
				http.Error(res, "unable to parse value as int", http.StatusBadRequest)
				return
			}
			storage.AddCounter(metric.Name, int64(val))
		default:
			http.Error(res, "unknown metric type, expected 'gauge' or 'counter'", http.StatusBadRequest)
			return
		}

		body := fmt.Sprintf("struct: %v+\n", storage)

		res.Write([]byte(body))
	}
}

func run() error {
	mux := http.NewServeMux()

	storage := MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}

	mux.HandleFunc("/", serveMetrics(&storage))

	return http.ListenAndServe("0.0.0.0:8080", mux)

}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
