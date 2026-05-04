package repository

import (
	"encoding/json"
	"errors"
	"os"
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
	return errors.New("wrong storage")
}

func (s *MemStorage) SaveToFile(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := FileData{
		Gauge:   s.gauge,
		Counter: s.counter,
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func (s *MemStorage) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var data FileData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range data.Gauge {
		s.gauge[k] = v
	}
	for k, v := range data.Counter {
		s.counter[k] = v
	}

	return nil
}

func (s *MemStorage) UpdateBatch(metrics []models.Update) error {
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

func NewMemStorage() *MemStorage {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}
