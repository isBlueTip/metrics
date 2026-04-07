package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestGetAllMetrics(t *testing.T) {
	type args struct {
		s      repository.Storage
		method string
	}
	type want struct {
		status      int
		response    string
		contentType string
	}
	tests := []struct {
		name  string
		args  args
		setup func(s repository.Storage)
		want  want
	}{
		{
			name: "1 positive",
			args: args{
				s:      repository.NewStorage(),
				method: http.MethodGet,
			},
			setup: func(s repository.Storage) {

			},
			want: want{
				status:      http.StatusOK,
				contentType: "text/html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.args.s)
			}
			router := chi.NewRouter()

			handler := GetAllMetrics(tt.args.s)
			router.Get("/", handler)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.args.method, "/", nil)
			router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			_, err := io.ReadAll(res.Body)

			require.Equal(t, tt.want.status, res.StatusCode, fmt.Sprintf("statusCode: %d, want %d\n", res.StatusCode, tt.want.status))

			gotCT := res.Header.Get("Content-Type")
			wantCT := tt.want.contentType
			if !strings.HasPrefix(gotCT, wantCT) {
				t.Errorf("contentType: %q, want starting with %q\n", gotCT, tt.want.contentType)
			}

			//if strings.HasPrefix(tt.want.contentType, "text/html") {
			//	//
			//} else if strings.HasPrefix(tt.want.contentType, "application/json") {
			//	assert.JSONEq(t, tt.want.response, string(resBody))
			//}

			require.NoError(t, err)

		})
	}
}

func TestGetMetricByName(t *testing.T) {
	type args struct {
		s      repository.Storage
		url    string
		method string
	}
	type want struct {
		status      int
		response    string
		contentType string
	}
	tests := []struct {
		name  string
		args  args
		setup func(s repository.Storage)
		want  want
	}{
		{
			name: "1 positive",
			args: args{
				s:      repository.NewStorage(),
				url:    "/value/gauge/metric1",
				method: http.MethodGet,
			},
			setup: func(s repository.Storage) {
				s.SetGauge("metric1", 158.3)
			},
			want: want{
				status:      http.StatusOK,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "2 positive",
			args: args{
				s:      repository.NewStorage(),
				url:    "/value/gauge/metric1",
				method: http.MethodGet,
			},
			want: want{
				status:      http.StatusNotFound,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "3 positive",
			args: args{
				s:      repository.NewStorage(),
				url:    "/value/counter/metric1",
				method: http.MethodGet,
			},
			setup: func(s repository.Storage) {
				s.SetCounter("metric1", 158)
			},
			want: want{
				status:      http.StatusOK,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "4 positive",
			args: args{
				s:      repository.NewStorage(),
				url:    "/value/counter/metric1",
				method: http.MethodGet,
			},
			want: want{
				status:      http.StatusNotFound,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "5 positive",
			args: args{
				s:      repository.NewStorage(),
				url:    "/value/unknown/metric1",
				method: http.MethodGet,
			},
			want: want{
				status:      http.StatusBadRequest,
				response:    "",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.args.s)
			}
			router := chi.NewRouter()

			handler := GetMetricByName(tt.args.s)
			router.Get("/value/{metricType}/{metricName}", handler)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			_, err := io.ReadAll(res.Body)

			require.Equal(t, tt.want.status, res.StatusCode, fmt.Sprintf("statusCode: %d, want %d\n", res.StatusCode, tt.want.status))

			gotCT := res.Header.Get("Content-Type")
			wantCT := tt.want.contentType
			if !strings.HasPrefix(gotCT, wantCT) {
				t.Errorf("contentType: %q, want starting with %q\n", gotCT, tt.want.contentType)
			}

			//if strings.HasPrefix(tt.want.contentType, "text/plain") {
			//	//
			//} else if strings.HasPrefix(tt.want.contentType, "application/json") {
			//	assert.JSONEq(t, tt.want.response, string(resBody))
			//}

			require.NoError(t, err)
		})
	}
}
