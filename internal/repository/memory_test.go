package repository

import (
	"testing"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/stretchr/testify/assert"
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
			name: "1 positive",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: models.Gauge,
				val:  15.6,
			},
		},
		{
			name: "2 positive empty name",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: "",
				val:  15.6,
			},
		},
		{
			name: "3 positive zero value",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: models.Gauge,
				val:  0,
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
			name: "1 positive",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: models.Counter,
				val:  15,
			},
		},
		{
			name: "2 positive empty name",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: "",
				val:  15,
			},
		},
		{
			name: "3 positive zero value",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				name: models.Counter,
				val:  0,
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
			assert.Equalf(t, tt.args.val, s.counter[tt.args.name], "MemStorage.counter[\"%s\"] = %v, want %v\n", tt.args.name, s.counter[tt.args.name], tt.args.val)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "1 positive",
			fields: fields{
				gauge:   map[string]float64{"metric_1": 15.8, "metric_2": -15.8},
				counter: map[string]int64{},
			},
			args: args{
				name: "metric_1",
			},
			want: 15.8,
		},
		{
			name: "2 positive",
			fields: fields{
				gauge:   map[string]float64{"metric_1": 15.8, "metric_2": -15.8},
				counter: map[string]int64{},
			},
			args: args{
				name: "metric_2",
			},
			want: -15.8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			if got, _ := s.GetGauge(tt.args.name); got != tt.want {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "1 positive",
			fields: fields{
				gauge:   map[string]float64{},
				counter: map[string]int64{"metric_1": 15, "metric_2": -15},
			},
			args: args{
				name: "metric_1",
			},
			want: 15,
		},
		{
			name: "2 positive",
			fields: fields{
				gauge:   map[string]float64{},
				counter: map[string]int64{"metric_1": 15, "metric_2": -15},
			},
			args: args{
				name: "metric_2",
			},
			want: -15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			if got, _ := s.GetCounter(tt.args.name); got != tt.want {
				t.Errorf("MemStorage.GetCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetGauges(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []models.GaugeModel
		wantErr bool
	}{
		{
			name: "1 positive",
			fields: fields{
				gauge: map[string]float64{
					"metric_1": 64.6,
					"metric_2": 5933.6515616,
				},
			},
			wantRes: []models.GaugeModel{
				{Name: "metric_1", Value: 64.6},
				{Name: "metric_2", Value: 5933.6515616},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			gotRes, err := s.GetGauges()
			if tt.wantErr {
				require.Error(t, err, "MemStorage.GetGauges() err: %v, wantErr %t", err, tt.wantErr)
			}
			assert.ElementsMatch(t, gotRes, tt.wantRes, "MemStorage.GetGauges() = %v, want %v", gotRes, tt.wantRes)
		})
	}
}

func TestMemStorage_GetCounters(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []models.CounterModel
		wantErr bool
	}{
		{
			name: "1 positive",
			fields: fields{
				counter: map[string]int64{
					"metric_3": 18,
					"metric_4": 5165151,
				},
			},
			wantRes: []models.CounterModel{
				{Name: "metric_3", Value: 18},
				{Name: "metric_4", Value: 5165151},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			gotRes, err := s.GetCounters()
			if tt.wantErr {
				require.Error(t, err, "MemStorage.GetCounters() err: %v, wantErr %t", err, tt.wantErr)
			}
			assert.ElementsMatch(t, gotRes, tt.wantRes, "MemStorage.GetCounters() = %v, want %v", gotRes, tt.wantRes)
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	want := &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
	got := NewMemStorage()
	require.EqualValuesf(t, got, want, "NewStorage() = %v, want %v", got, want)

}
