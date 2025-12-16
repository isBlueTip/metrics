package main

type StorageInterface interface {
	AddGauge(name string, val float64)
	AddCounter(name string, val int64)
}

type MemStorage struct {
	StorageInterface
	Gauge   map[string]float64
	Counter map[string]int64
}

func (s *MemStorage) AddGauge(name string, val float64) {
	s.Gauge[name] = val
}

func (s *MemStorage) AddCounter(name string, val int64) {
	s.Counter[name] += val
}
