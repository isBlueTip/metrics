package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io"
	"net/http/httptest"
	"testing"

	"github.com/isBlueTip/metrics/internal/repository"
)

func TestUpdateMetric(t *testing.T) {
	storage := repository.NewStorage()
	type args struct {
		url    string
		method string
	}
	type want struct {
		status      int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				url:    "/update/gauge/Alloc/9",
				method: http.MethodPost,
			},
			want: want{
				status:      http.StatusOK,
				response:    "{}",
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()

			handler := UpdateMetric(storage)
			router.Post("/update/{metricType}/{metricName}/{metricVal}", handler)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.Equal(t, tt.want.status, res.StatusCode, fmt.Sprintf("body: %s\n", resBody))

			require.NoError(t, err)

			assert.JSONEq(t, tt.want.response, string(resBody))
			//assert.Equal(t, tt.want.contentType, res.ContentType)

		})
	}
}
