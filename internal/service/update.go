package service

import (
	"github.com/xxx/metrics/internal/repository"
)

func UpdateCounter(storage repository.StorageInterface, name string, val int64) {
	storage.AddCounter(name, val)
}

func UpdateGauge(storage repository.StorageInterface, name string, val float64) {
	storage.AddGauge(name, val)
}
