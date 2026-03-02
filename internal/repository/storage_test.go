package repository

import (
	"testing"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_SetGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name string
		val  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: models.Gauge,
				val:  15.6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			s.SetGauge(tt.args.name, tt.args.val)
			require.Equalf(t, tt.args.val, s.gauge[tt.args.name], "MemStorage.gauge[\"%s\"] = %v, want %v\n", tt.args.name, s.gauge[tt.args.name], tt.args.val)
		})
	}
}

func TestMemStorage_SetCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name string
		val  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: models.Counter,
				val:  15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			s.SetCounter(tt.args.name, tt.args.val)
			require.Equalf(t, tt.args.val, s.counter[tt.args.name], "MemStorage.counter[\"%s\"] = %v, want %v\n", tt.args.name, s.counter[tt.args.name], tt.args.val)
		})
	}
}

func TestNewStorage(t *testing.T) {
	want := &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
	got := NewStorage()
	require.EqualValuesf(t, got, want, "NewStorage() = %v, want %v", got, want)

}
