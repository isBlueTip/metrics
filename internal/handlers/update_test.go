package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"io"
	"net/http/httptest"
	"testing"

	"github.com/isBlueTip/metrics/internal/repository"
)

func TestUpdateMetric(t *testing.T) {
	storage := repository.MemStorage{}
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
				url:    "/update/gauge/r/9/",
				method: "POST",
			},
			want: want{
				status:      200,
				response:    "",
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()

			handler := UpdateMetric(&storage)
			handler(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.status, res.StatusCode)
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			assert.JSONEq(t, tt.want.response, string(resBody))
			//assert.Equal(t, tt.want.contentType, res.ContentType)
		})
	}
}
