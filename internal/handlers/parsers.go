package handlers

import (
	"fmt"
	"strconv"

	"github.com/isBlueTip/metrics/internal/models"
)

// todo get rid of interface{}
func ParseValue(metricType string, value string) (interface{}, error) {
	switch metricType {
	case models.Gauge:
		return strconv.ParseFloat(value, 64)
	case models.Counter:
		return strconv.ParseInt(value, 10, 64)
	default:
		return nil, fmt.Errorf("unknown metric type: %s, expected '%s' or '%s'", metricType, models.Gauge, models.Counter)
	}
}
