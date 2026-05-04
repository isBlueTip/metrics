package repository

import (
	"sync"

	"github.com/isBlueTip/metrics/internal/models"
)

type SyncStorage struct {
	mu       sync.RWMutex
	storage  Storage
	filePath string
}

func (s *SyncStorage) SetGauge(name string, val float64) error {
	if err := s.storage.SetGauge(name, val); err != nil {
		return err
	}
	return s.storage.SaveToFile(s.filePath)
}

func (s *SyncStorage) SetCounter(name string, val int64) error {
	if err := s.storage.SetCounter(name, val); err != nil {
		return err
	}
	return s.storage.SaveToFile(s.filePath)
}

func (s *SyncStorage) GetGauge(name string) (val float64, exists bool) {
	return s.storage.GetGauge(name)
}

func (s *SyncStorage) GetCounter(name string) (val int64, exists bool) {
	return s.storage.GetCounter(name)
}

func (s *SyncStorage) GetGauges() (res []models.GaugeModel, err error) {
	return s.storage.GetGauges()
}

func (s *SyncStorage) GetCounters() (res []models.CounterModel, err error) {
	return s.storage.GetCounters()
}

func (s *SyncStorage) Ping() error {
	return s.storage.Ping()
}

func (s *SyncStorage) SaveToFile(path string) error {
	return s.storage.SaveToFile(path)
}

func (s *SyncStorage) LoadFromFile(path string) error {
	return s.storage.LoadFromFile(path)
}

func (s *SyncStorage) UpdateBatch(metrics []models.Update) error {
	var err error
	for _, metric := range metrics {
		if metric.MType == models.Gauge {
			err = s.SetGauge(metric.ID, *metric.Value)

		} else if metric.MType == models.Counter {
			err = s.SetCounter(metric.ID, *metric.Delta)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func NewSyncStorage(storage Storage, filePath string) *SyncStorage {
	return &SyncStorage{storage: storage, filePath: filePath}
}
