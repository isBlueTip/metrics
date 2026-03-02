package service

import (
	"testing"

	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestUpdateGauge(t *testing.T) {
	type args struct {
		storage repository.Storage
		name    string
		val     float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1 positive",
			args: args{
				storage: repository.NewStorage(),
				name:    "metric_1",
				val:     125.89,
			},
		},
		{
			name: "2 positive zero",
			args: args{
				storage: repository.NewStorage(),
				name:    "metric_1",
				val:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateGauge(tt.args.storage, tt.args.name, tt.args.val)
			got := tt.args.storage.GetGauge(tt.args.name)
			assert.Equalf(t, tt.args.val, got, "MemStorage.GetGauge(): %v, want %v\n", got, tt.args.val)
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	type args struct {
		storage repository.Storage
		name    string
		val     int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1 positive",
			args: args{
				storage: repository.NewStorage(),
				name:    "metric_1",
				val:     512,
			},
		},
		{
			name: "2 positive zero",
			args: args{
				storage: repository.NewStorage(),
				name:    "metric_1",
				val:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateCounter(tt.args.storage, tt.args.name, tt.args.val)
			got := tt.args.storage.GetCounter(tt.args.name)
			assert.Equalf(t, tt.args.val, got, "MemStorage.GetCounter(): %v, want %v\n", got, tt.args.val)
		})
	}
}
