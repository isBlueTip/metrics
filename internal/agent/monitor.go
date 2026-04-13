package agent

import (
	"math/rand"
	"runtime"
)

type MetricSet struct {
	//MemStat runtime.MemStats
	Uints       map[string]uint64
	Floats      map[string]float64
	RandomValue int64
}

func (m *MetricSet) Collect() {
	memstat := &runtime.MemStats{}
	runtime.ReadMemStats(memstat)

	m.Uints["Alloc"] = memstat.Alloc
	m.Uints["BuckHashSys"] = memstat.BuckHashSys
	m.Uints["Frees"] = memstat.Frees
	m.Floats["GCCPUFraction"] = memstat.GCCPUFraction
	m.Uints["GCSys"] = memstat.GCSys
	m.Uints["HeapAlloc"] = memstat.HeapAlloc
	m.Uints["HeapIdle"] = memstat.HeapIdle
	m.Uints["HeapInuse"] = memstat.HeapInuse
	m.Uints["HeapObjects"] = memstat.HeapObjects
	m.Uints["HeapReleased"] = memstat.HeapReleased
	m.Uints["HeapSys"] = memstat.HeapSys
	m.Uints["LastGC"] = memstat.LastGC
	m.Uints["Lookups"] = memstat.Lookups
	m.Uints["MCacheInuse"] = memstat.MCacheInuse
	m.Uints["MCacheSys"] = memstat.MCacheSys
	m.Uints["MSpanInuse"] = memstat.MSpanInuse
	m.Uints["MSpanSys"] = memstat.MSpanSys
	m.Uints["Mallocs"] = memstat.Mallocs
	m.Uints["NextGC"] = memstat.NextGC
	m.Uints["NumForcedGC"] = uint64(memstat.NumForcedGC)
	m.Uints["NumGC"] = uint64(memstat.NumGC)
	m.Uints["OtherSys"] = memstat.OtherSys
	m.Uints["PauseTotalNs"] = memstat.PauseTotalNs
	m.Uints["StackInuse"] = memstat.StackInuse
	m.Uints["StackSys"] = memstat.StackSys
	m.Uints["Sys"] = memstat.Sys
	m.Uints["TotalAlloc"] = memstat.TotalAlloc
	m.Uints["PollCount"] += 1

	m.RandomValue = rand.Int63()
}
