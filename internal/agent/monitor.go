package agent

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

type MetricSet struct {
	//runtime.MemStats
	Alloc         uint64
	BuckHashSys   uint64
	Frees         uint64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     uint64
	HeapIdle      uint64
	HeapInuse     uint64
	HeapObjects   uint64
	HeapReleased  uint64
	HeapSys       uint64
	LastGC        uint64
	Lookups       uint64
	MCacheInuse   uint64
	MCacheSys     uint64
	MSpanInuse    uint64
	MSpanSys      uint64
	Mallocs       uint64
	NextGC        uint64
	NumForcedGC   uint32
	NumGC         uint32
	OtherSys      uint64
	PauseTotalNs  uint64
	StackInuse    uint64
	StackSys      uint64
	Sys           uint64
	TotalAlloc    uint64
	PollCount     int
	RandomValue   int64
}

func (m *MetricSet) Collect() {
	memstat := &runtime.MemStats{}
	runtime.ReadMemStats(memstat)

	m.Alloc = memstat.Alloc
	m.BuckHashSys = memstat.BuckHashSys
	m.Frees = memstat.Frees
	m.GCCPUFraction = memstat.GCCPUFraction
	m.GCSys = memstat.GCCPUFraction
	m.HeapAlloc = memstat.HeapAlloc
	m.HeapIdle = memstat.HeapIdle
	m.HeapInuse = memstat.HeapInuse
	m.HeapObjects = memstat.HeapObjects
	m.HeapReleased = memstat.HeapReleased
	m.HeapSys = memstat.HeapSys
	m.LastGC = memstat.LastGC
	m.Lookups = memstat.Lookups
	m.MCacheInuse = memstat.MCacheInuse
	m.MCacheSys = memstat.MCacheSys
	m.MSpanInuse = memstat.MSpanInuse
	m.MSpanSys = memstat.MSpanSys
	m.Mallocs = memstat.Mallocs
	m.NextGC = memstat.NextGC
	m.NumForcedGC = memstat.NumForcedGC
	m.NumGC = memstat.NumGC
	m.OtherSys = memstat.OtherSys
	m.PauseTotalNs = memstat.PauseTotalNs
	m.StackInuse = memstat.StackInuse
	m.StackSys = memstat.StackSys
	m.Sys = memstat.Sys
	m.TotalAlloc = memstat.TotalAlloc

	x := rand.New(rand.NewSource(time.Now().UnixNano()))
	m.PollCount += 1
	m.RandomValue = x.Int63()

	log.Printf("metricSet collected: %+v\n", *m)
}
