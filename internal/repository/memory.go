package repository

import (
	"sync"

	"github.com/isBlueTip/metrics/internal/models"
)

type MemStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func (s *MemStorage) SetGauge(name string, val float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauge[name] = val
	return nil
}

func (s *MemStorage) SetCounter(name string, val int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter[name] += val
	return nil
}

func (s *MemStorage) GetGauge(name string) (val float64, exists bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists = s.gauge[name]
	return
}

func (s *MemStorage) GetCounter(name string) (val int64, exists bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists = s.counter[name]
	return
}

func (s *MemStorage) GetGauges() (res []models.GaugeModel, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for n, v := range s.gauge {
		record := models.GaugeModel{Name: n, Value: v}
		res = append(res, record)
	}

	return
}

func (s *MemStorage) GetCounters() (res []models.CounterModel, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for n, v := range s.counter {
		record := models.CounterModel{Name: n, Value: v}
		res = append(res, record)
	}

	return
}

func (s *MemStorage) Ping() error {
	return nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}
