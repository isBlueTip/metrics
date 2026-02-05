package agent

import "testing"

func TestMetricSet_Collect(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		//{
		//	name: "1 positive",
		//	fields: fields{
		//		Alloc:         1885165,
		//		BuckHashSys:   56165165,
		//		Frees:         498464,
		//		GCCPUFraction: 6851681.849861,
		//		GCSys:         615651.51564165,
		//		HeapAlloc:     651651,
		//		HeapIdle:      65416168,
		//		HeapInuse:     65165,
		//		HeapObjects:   561651,
		//		HeapReleased:  1561651,
		//		HeapSys:       84894864,
		//		LastGC:        846651,
		//		Lookups:       8468651,
		//		MCacheInuse:   46846,
		//		MCacheSys:     84651,
		//		MSpanInuse:    6541651,
		//		MSpanSys:      86546,
		//		Mallocs:       23354,
		//		NextGC:        651,
		//		NumForcedGC:   31556,
		//		NumGC:         1313216,
		//		OtherSys:      516516,
		//		PauseTotalNs:  1561321,
		//		StackInuse:    3545451,
		//		StackSys:      1354841,
		//		Sys:           1056886,
		//		TotalAlloc:    86415,
		//		PollCount:     5,
		//		RandomValue:   1641352164,
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricSet{
				Alloc:         tt.fields.Alloc,
				BuckHashSys:   tt.fields.BuckHashSys,
				Frees:         tt.fields.Frees,
				GCCPUFraction: tt.fields.GCCPUFraction,
				GCSys:         tt.fields.GCSys,
				HeapAlloc:     tt.fields.HeapAlloc,
				HeapIdle:      tt.fields.HeapIdle,
				HeapInuse:     tt.fields.HeapInuse,
				HeapObjects:   tt.fields.HeapObjects,
				HeapReleased:  tt.fields.HeapReleased,
				HeapSys:       tt.fields.HeapSys,
				LastGC:        tt.fields.LastGC,
				Lookups:       tt.fields.Lookups,
				MCacheInuse:   tt.fields.MCacheInuse,
				MCacheSys:     tt.fields.MCacheSys,
				MSpanInuse:    tt.fields.MSpanInuse,
				MSpanSys:      tt.fields.MSpanSys,
				Mallocs:       tt.fields.Mallocs,
				NextGC:        tt.fields.NextGC,
				NumForcedGC:   tt.fields.NumForcedGC,
				NumGC:         tt.fields.NumGC,
				OtherSys:      tt.fields.OtherSys,
				PauseTotalNs:  tt.fields.PauseTotalNs,
				StackInuse:    tt.fields.StackInuse,
				StackSys:      tt.fields.StackSys,
				Sys:           tt.fields.Sys,
				TotalAlloc:    tt.fields.TotalAlloc,
				PollCount:     tt.fields.PollCount,
				RandomValue:   tt.fields.RandomValue,
			}
			m.Collect()
		})
	}
}
