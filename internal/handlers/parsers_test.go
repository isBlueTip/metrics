package handlers

import (
	"reflect"
	"testing"

	"github.com/isBlueTip/metrics/internal/models"
)

func TestParseValue(t *testing.T) {
	type args struct {
		metricType string
		value      string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "positive Gauge",
			args: args{
				metricType: models.Gauge,
				value:      "15.5",
			},
			want:    15.5,
			wantErr: false,
		},
		{
			name: "positive Counter",
			args: args{
				metricType: models.Counter,
				value:      "16",
			},
			want:    int64(16),
			wantErr: false,
		},
		{
			name: "negative Unknown",
			args: args{
				metricType: "Unknown",
				value:      "16",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValue(tt.args.metricType, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
