package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Gauge(t *testing.T) {
	type args struct {
		storage    repository.Storage
		path       string
		metricName string
	}
	type want struct {
		value      float64
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 positive",
			args: args{
				storage:    repository.NewStorage(),
				path:       "/update/" + models.Gauge + "/metric_1" + "/101.1",
				metricName: "metric_1",
			},
			want: want{
				value:      101.1,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "2 negative",
			args: args{
				storage:    repository.NewStorage(),
				path:       "/update/" + "invalidType" + "/metric_1" + "/101.1",
				metricName: "metric_1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "3 negative",
			args: args{
				storage:    repository.NewStorage(),
				path:       "/delete",
				metricName: "metric_1",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := Router(tt.args.storage, "")

			req := httptest.NewRequest(http.MethodPost, tt.args.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			got, _ := tt.args.storage.GetGauge(tt.args.metricName)
			assert.Equalf(t, tt.want.statusCode, w.Code, "POST /update/%s/%s/%v/ err: statusCode = %v, want %v\n", models.Gauge, tt.args.metricName, tt.want.value, w.Code, tt.want.statusCode)

			if tt.want.statusCode == http.StatusOK {
				assert.Equalf(t, tt.want.value, got, "POST /update/%s/%s/%v/ err: storage.GetGauge(%s) = %v, want %v\n", models.Gauge, tt.args.metricName, tt.want.value, tt.args.metricName, got, tt.want.value)
			}

		})
	}
}

func TestRouter_Counter(t *testing.T) {
	type args struct {
		storage    repository.Storage
		path       string
		metricName string
	}
	type want struct {
		value      int64
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 positive",
			args: args{
				storage:    repository.NewStorage(),
				path:       "/update/" + models.Counter + "/metric_1" + "/101",
				metricName: "metric_1",
			},
			want: want{
				value:      101,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "2 negative",
			args: args{
				storage:    repository.NewStorage(),
				path:       "/update" + "/invalidType" + "/metric_1" + "/101",
				metricName: "metric_1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "3 negative",
			args: args{
				storage:    repository.NewStorage(),
				path:       "/delete",
				metricName: "metric_1",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Router(tt.args.storage, "")

			req := httptest.NewRequest(http.MethodPost, tt.args.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			got, _ := tt.args.storage.GetCounter(tt.args.metricName)
			assert.Equalf(t, tt.want.statusCode, w.Code, "POST %s err: statusCode = %v, want %v\n", tt.args.path, w.Code, tt.want.statusCode)

			if tt.want.statusCode == http.StatusOK {
				assert.Equalf(t, tt.want.value, got, "POST %s err: storage.GetCounter(%s) = %v, want %v\n", tt.args.path, tt.args.metricName, got, tt.want.value)
			}
			t.Logf("tt.args.storage = %+v\n", tt.args.storage)

		})
	}
}
