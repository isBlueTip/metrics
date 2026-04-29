package repository

import (
	"github.com/isBlueTip/metrics/internal/models"
)

type Storage interface {
	SetGauge(name string, val float64) error
	SetCounter(name string, val int64) error
	GetGauge(name string) (val float64, exists bool)
	GetCounter(name string) (val int64, exists bool)
	GetGauges() (val []models.GaugeModel, err error)
	GetCounters() (val []models.CounterModel, err error)
	Ping() error
	//SaveToFile(path string) error
	//LoadFromFile(path string) error
}
