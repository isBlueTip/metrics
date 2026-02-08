package agent

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/isBlueTip/metrics/internal/models"
)

const serverAddr = "127.0.0.1:8080"

func TestSender_Send(t *testing.T) {
	type fields struct {
		HC   http.Client
		Addr string
	}
	type args struct {
		metricSet *MetricSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1 positive",
			fields: fields{
				HC:   http.Client{},
				Addr: serverAddr,
			},
			args: args{
				metricSet: &MetricSet{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sender{
				HC:   tt.fields.HC,
				Addr: serverAddr,
			}
			if err := s.Send(tt.args.metricSet); (err != nil) != tt.wantErr {
				t.Errorf("Sender.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateURL(t *testing.T) {
	type args struct {
		metricType string
		metricName string
		metricVal  string
	}
	tests := []struct {
		name string
		args args
		want *url.URL
	}{
		{name: "1 positive gauge",
			args: args{
				metricType: models.Gauge,
				metricName: "testName1",
				metricVal:  "8558",
			},
			want: &url.URL{
				Scheme: "http",
				Host:   serverAddr,
				Path:   "/update/" + models.Gauge + "/" + "testName1" + "/" + "8558",
			},
		},
		{name: "2 positive counter",
			args: args{
				metricType: models.Counter,
				metricName: "testName2",
				metricVal:  "8560",
			},
			want: &url.URL{
				Scheme: "http",
				Host:   serverAddr,
				Path:   "/update/" + models.Counter + "/" + "testName2" + "/" + "8560",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateURL(serverAddr, tt.args.metricType, tt.args.metricName, tt.args.metricVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executeRequest(t *testing.T) {
	sender := Sender{
		HC:   http.Client{},
		Addr: serverAddr,
	}
	type args struct {
		u  *url.URL
		hc http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "status 200",
			args: args{
				u: &url.URL{
					Scheme: "http",
					Host:   serverAddr,
					Path:   "/update/" + models.Counter + "/" + "testName2" + "/" + "8560",
				},
				hc: http.Client{},
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "status 404",
			args: args{
				u: &url.URL{
					Scheme: "http",
					Host:   serverAddr,
					Path:   "/update/" + models.Counter + "/" + "" + "/" + "8560",
				},
				hc: http.Client{},
			},
			wantErr: false,
		},
		{
			name: "status 400",
			args: args{
				u: &url.URL{
					Scheme: "http",
					Host:   serverAddr,
					Path:   "/update/" + "models.Counter" + "/" + "testName2" + "/" + "8560",
				},
				hc: http.Client{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := executeRequest(tt.args.u, sender); (err != nil) != tt.wantErr {
				t.Errorf("executeRequest() error = %v, wantError %v", err, tt.wantErr)
			}
		})
	}
}
