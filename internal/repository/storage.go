package repository

import (
	"sync"

	"github.com/isBlueTip/metrics/internal/models"
)

type Storage interface {
	SetGauge(name string, val float64)
	SetCounter(name string, val int64)
	GetGauge(name string) (val float64, exists bool)
	GetCounter(name string) (val int64, exists bool)
	GetGauges() (val []models.GaugeModel)
	GetCounters() (val []models.CounterModel)
}

type MemStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func (s *MemStorage) SetGauge(name string, val float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauge[name] = val
}

func (s *MemStorage) SetCounter(name string, val int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter[name] += val
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

func (s *MemStorage) GetGauges() (res []models.GaugeModel) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for n, v := range s.gauge {
		record := models.GaugeModel{Name: n, Value: v}
		res = append(res, record)
	}
	return
}

func (s *MemStorage) GetCounters() (res []models.CounterModel) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for n, v := range s.counter {
		record := models.CounterModel{Name: n, Value: v}
		res = append(res, record)
	}
	return
}

func NewStorage() *MemStorage {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}
