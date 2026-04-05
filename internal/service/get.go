package service

import (
	"fmt"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
)

func GetGauge(storage repository.Storage, name string) (val float64, err error) {
	val, exists := storage.GetGauge(name)
	if !exists {
		err = fmt.Errorf("value not found: %s", name)
	}

	return
}

func GetCounter(storage repository.Storage, name string) (val int64, err error) {
	val, exists := storage.GetCounter(name)
	if !exists {
		err = fmt.Errorf("value not found: %s", name)
	}
	return
}

func GetGauges(storage repository.Storage) []models.GaugeModel {
	return storage.GetGauges()
}

func GetCounters(storage repository.Storage) []models.CounterModel {
	return storage.GetCounters()
}
