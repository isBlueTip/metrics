package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"io"
	"net/http/httptest"
	"testing"

	"github.com/isBlueTip/metrics/internal/repository"
)

func TestUpdate(t *testing.T) {
	storage := repository.NewMemStorage()
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
			name: "1 positive Gauge",
			args: args{
				url:    "/update/gauge/Alloc/9",
				method: http.MethodPost,
			},
			want: want{
				status:   http.StatusOK,
				response: "{}",
				//contentType: "application/json",
				contentType: "",
			},
		},
		{
			name: "2 positive Counter",
			args: args{
				url:    "/update/counter/PollCount/9",
				method: http.MethodPost,
			},
			want: want{
				status:   http.StatusOK,
				response: "{}",
				//contentType: "application/json",
				contentType: "",
			},
		},
		{
			name: "3 negative unknown metricType",
			args: args{
				url:    "/update/unknownType/Alloc/9",
				method: http.MethodPost,
			},
			want: want{
				status: http.StatusBadRequest,
				//response:    "{}",
				//contentType: "application/json",
				contentType: "",
			},
		},
		{
			name: "4 negative invalid metricVal",
			args: args{
				url:    "/update/gauge/Alloc/abc",
				method: http.MethodPost,
			},
			want: want{
				status: http.StatusBadRequest,
				//response:    "{}",
				//contentType: "application/json",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()

			handler := UpdateURL(storage)
			router.Post("/update/{metricType}/{metricName}/{metricVal}", handler)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			router.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			_, err := io.ReadAll(res.Body)

			require.Equal(t, tt.want.status, res.StatusCode, fmt.Sprintf("statusCode: %d, want %d\n", res.StatusCode, tt.want.status))

			gotCT := res.Header.Get("ContentType")
			wantCT := tt.want.contentType
			if !strings.HasPrefix(gotCT, wantCT) {
				t.Errorf("contentType: %q, want starting with %q\n", gotCT, tt.want.contentType)
			}

			require.NoError(t, err)

			//assert.JSONEq(t, tt.want.response, string(resBody))

		})
	}
}
