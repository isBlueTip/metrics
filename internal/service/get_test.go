package service

import (
	"testing"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGauge(t *testing.T) {
	type args struct {
		storage repository.Storage
		name    string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(s repository.Storage)
		wantVal float64
		wantErr bool
	}{
		{
			name: "1 positive",
			args: args{
				storage: repository.NewMemStorage(),
				name:    "metric_1",
			},
			setup: func(s repository.Storage) {
				s.SetGauge("metric_1", 159.652)
			},
			wantVal: 159.652,
			wantErr: false,
		},
		{
			name: "2 negative",
			args: args{
				storage: repository.NewMemStorage(),
				name:    "metric_1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.args.storage)
			}
			gotVal, err := GetGauge(tt.args.storage, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("GetGauge() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestGetCounter(t *testing.T) {
	type args struct {
		storage repository.Storage
		name    string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(s repository.Storage)
		wantVal int64
		wantErr bool
	}{
		{
			name: "1 positive",
			args: args{
				storage: repository.NewMemStorage(),
				name:    "metric_1",
			},
			setup: func(s repository.Storage) {
				s.SetCounter("metric_1", 65236)
			},
			wantVal: 65236,
			wantErr: false,
		},
		{
			name: "2 negative",
			args: args{
				storage: repository.NewMemStorage(),
				name:    "metric_1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.args.storage)
			}
			gotVal, err := GetCounter(tt.args.storage, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("GetCounter() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestGetGauges(t *testing.T) {
	type args struct {
		storage repository.Storage
	}
	tests := []struct {
		name    string
		args    args
		setup   func(s repository.Storage)
		want    []models.GaugeModel
		wantErr bool
	}{
		{
			name: "1 positive",
			args: args{
				storage: repository.NewMemStorage(),
			},
			setup: func(s repository.Storage) {
				s.SetGauge("metric1", 15.2)
				s.SetGauge("metric2", 155.2)
				s.SetGauge("metric3", 1569.2)
				s.SetGauge("metric4", 155683.2)
			},
			want: []models.GaugeModel{
				{Name: "metric1", Value: 15.2},
				{Name: "metric2", Value: 155.2},
				{Name: "metric3", Value: 1569.2},
				{Name: "metric4", Value: 155683.2},
			},
			wantErr: false,
		},
		{
			name: "2 positive",
			args: args{
				storage: repository.NewMemStorage(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.args.storage)
			}

			got, err := GetGauges(tt.args.storage)
			if tt.wantErr {
				require.Error(t, err, "GetGauges() err: %v, wantErr %t", err, tt.wantErr)
			}
			assert.ElementsMatch(t, got, tt.want, "GetGauges() = %v, want %v", got, tt.want)
		})
	}
}

func TestGetCounters(t *testing.T) {
	type args struct {
		storage repository.Storage
	}
	tests := []struct {
		name    string
		args    args
		setup   func(s repository.Storage)
		want    []models.CounterModel
		wantErr bool
	}{
		{
			name: "1 positive",
			args: args{
				storage: repository.NewMemStorage(),
			},
			setup: func(s repository.Storage) {
				s.SetCounter("metric1", 15)
				s.SetCounter("metric2", 155)
				s.SetCounter("metric3", 1569)
				s.SetCounter("metric4", 155683)

			},
			want: []models.CounterModel{
				{Name: "metric1", Value: 15},
				{Name: "metric2", Value: 155},
				{Name: "metric3", Value: 1569},
				{Name: "metric4", Value: 155683},
			},
			wantErr: false,
		},
		{
			name: "2 positive",
			args: args{
				storage: repository.NewMemStorage(),
			},
			setup: func(s repository.Storage) {

			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.args.storage)
			}

			got, err := GetCounters(tt.args.storage)
			if tt.wantErr {
				require.Error(t, err, "GetCounters() err: %v, wantErr %t", err, tt.wantErr)
			}
			assert.ElementsMatch(t, got, tt.want, "GetCounters() = %v, want %v", got, tt.want)
		})
	}
}
