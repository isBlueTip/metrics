package repository

type Storage interface {
	SetGauge(name string, val float64)
	SetCounter(name string, val int64)
	GetGauge(name string) (val float64)
	GetCounter(name string) (val int64)
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

func (s *MemStorage) GetGauge(name string) float64 {
	return s.gauge[name]
}

func (s *MemStorage) GetCounter(name string) int64 {
	return s.counter[name]
}

func NewStorage() Storage {
	return &MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)}
}
