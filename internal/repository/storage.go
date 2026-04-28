package repository

import (
	"github.com/isBlueTip/metrics/internal/models"
)

type Storage interface {
	SetGauge(name string, val float64)
	SetCounter(name string, val int64)
	GetGauge(name string) (val float64, exists bool)
	GetCounter(name string) (val int64, exists bool)
	GetGauges() (val []models.GaugeModel)
	GetCounters() (val []models.CounterModel)
	SaveToFile(path string) error
	LoadFromFile(path string) error
}
