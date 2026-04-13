package service

import (
	"github.com/isBlueTip/metrics/internal/repository"
)

func UpdateGauge(storage repository.Storage, name string, val float64) {
	storage.SetGauge(name, val)
}

func UpdateCounter(storage repository.Storage, name string, val int64) {
	storage.SetCounter(name, val)
}
