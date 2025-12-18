package handler

import (
	"errors"
	"fmt"
	"github.com/xxx/metrics/internal/models"
	"github.com/xxx/metrics/internal/repository"
	"github.com/xxx/metrics/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type URLMetric struct {
	Type, Name string
	Value      interface{}
}

func UpdateMetric(storage repository.StorageInterface) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Println("   **********************   ")
		log.Println("")
		if req.Method != http.MethodPost {
			http.Error(res, "only POST request accepted", http.StatusMethodNotAllowed)
			return
		}

		log.Printf("'%v', type: %T\r\n", req.URL.Path, req.URL.Path)
		splitURL := strings.Split(req.URL.Path, "/")[1:]

		// validate
		if len(splitURL) < 3 {
			http.Error(res, "not enough URL params to write metric, expected 3", http.StatusNotFound)
			return
		}
		parsedValue, err := ParseValue(splitURL[0], splitURL[2])
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				http.Error(res, "invalid metric value, expected numeric", http.StatusBadRequest)
			} else {
				http.Error(res, err.Error(), http.StatusBadRequest)
			}
			return
		}

		switch splitURL[0] {
		case models.Gauge:
			service.UpdateGauge(storage, splitURL[1], parsedValue.(float64))
		case models.Counter:
			service.UpdateCounter(storage, splitURL[1], parsedValue.(int64))
		}

		body := fmt.Sprintf("struct: %+v\n", storage)

		res.Write([]byte(body))
	}
}
