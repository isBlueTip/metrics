package service

import (
	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
)

func UpdateGauge(storage repository.Storage, name string, val float64) error {
	return storage.SetGauge(name, val)
}

func UpdateCounter(storage repository.Storage, name string, val int64) error {
	return storage.SetCounter(name, val)
}

func UpdateBatch(storage repository.Storage, metrics []models.Update) error {
	return storage.UpdateBatch(metrics)
}
