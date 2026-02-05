package handlers

import (
	"fmt"
	"github.com/xxx/metrics/internal/models"
	"strconv"
)

func ParseValue(metricType string, value string) (interface{}, error) {
	switch metricType {
	case models.Gauge:
		return strconv.ParseFloat(value, 64)
	case models.Counter:
		return strconv.ParseInt(value, 10, 64)
	default:
		return nil, fmt.Errorf("unknown metric type, expected '%s' or '%s'", models.Gauge, models.Counter)
	}
}
