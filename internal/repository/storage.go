package repository

import "github.com/isBlueTip/metrics/internal/models"

type Storage interface {
	SetGauge(name string, val float64)
	SetCounter(name string, val int64)
	GetGauge(name string) (val float64, exists bool)
	GetCounter(name string) (val int64, exists bool)
	GetGauges() (val []models.GaugeModel)
	GetCounters() (val []models.CounterModel)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (s *MemStorage) SetGauge(name string, val float64) {
	s.gauge[name] = val
}

func (s *MemStorage) SetCounter(name string, val int64) {
	s.counter[name] += val
}

func (s *MemStorage) GetGauge(name string) (val float64, exists bool) {
	val, exists = s.gauge[name]
	return
}

func (s *MemStorage) GetCounter(name string) (val int64, exists bool) {
	val, exists = s.counter[name]
	return
}

func (s *MemStorage) GetGauges() (res []models.GaugeModel) {
	for n, v := range s.gauge {
		record := models.GaugeModel{Name: n, Value: v}
		res = append(res, record)
	}
	return
}

func (s *MemStorage) GetCounters() (res []models.CounterModel) {
	for n, v := range s.counter {
		record := models.CounterModel{Name: n, Value: v}
		res = append(res, record)
	}
	return
}

func NewStorage() *MemStorage {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}
