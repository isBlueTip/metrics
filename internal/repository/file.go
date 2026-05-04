package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/isBlueTip/metrics/internal/models"
)

type FileStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
	path    string
}

type FileData struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
}

func (s *FileStorage) SetGauge(name string, val float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauge[name] = val
	return nil
}

func (s *FileStorage) SetCounter(name string, val int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter[name] += val
	return nil
}

func (s *FileStorage) GetGauge(name string) (val float64, exists bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists = s.gauge[name]
	return
}

func (s *FileStorage) GetCounter(name string) (val int64, exists bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists = s.counter[name]
	return
}

func (s *FileStorage) GetGauges() (res []models.GaugeModel, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for n, v := range s.gauge {
		record := models.GaugeModel{Name: n, Value: v}
		res = append(res, record)
	}
	return
}

func (s *FileStorage) GetCounters() (res []models.CounterModel, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for n, v := range s.counter {
		record := models.CounterModel{Name: n, Value: v}
		res = append(res, record)
	}
	return
}

func (s *FileStorage) Ping() error {
	return errors.New("wrong storage")
}

func (s *FileStorage) SaveToFile(path string) error {
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

func (s *FileStorage) LoadFromFile(path string) error {
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

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{gauge: make(map[string]float64), counter: make(map[string]int64), path: path}
}
