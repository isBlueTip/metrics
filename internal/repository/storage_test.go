package repository

import (
	"reflect"
	"testing"

	"github.com/isBlueTip/metrics/internal/models"
)

func TestMemStorage_AddGauge(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
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
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
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
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.AddGauge(tt.args.name, tt.args.val)
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
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
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
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
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			s.AddCounter(tt.args.name, tt.args.val)
		})
	}
}

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
