package repository

type Storage interface {
	SetGauge(name string, val float64)
	SetCounter(name string, val int64)
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

func NewStorage() Storage {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}
