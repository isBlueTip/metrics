package agent

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/isBlueTip/metrics/internal/models"
	"github.com/sethgrid/pester"
)

type Sender struct {
	HC   *pester.Client
	Addr string
}

func (s *Sender) executeRequest(u *url.URL) (*http.Response, error) {
	log.Printf("sending metric to url: %s\n", u.String())
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.HC.Do(req)
	if err != nil {
		return nil, err
	}

	log.Printf("status: %s\n", resp.Status)

	return resp, err
}

func (s *Sender) Send(metricSet *MetricSet) error {
	var u *url.URL
	var resp *http.Response
	var err error

	for k, v := range metricSet.Uints {
		switch k {
		case "PollCount":
			u, err = generateURL(s.Addr, models.Counter, k, strconv.FormatUint(v, 10))
		default:
			u, err = generateURL(s.Addr, models.Gauge, k, strconv.FormatUint(v, 10))
		}
		if err != nil {
			return err
		}
		resp, err = s.executeRequest(u)
		if err != nil {
			return err
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	for k, v := range metricSet.Floats {
		u, err = generateURL(s.Addr, models.Gauge, k, strconv.FormatFloat(v, 'f', 4, 64))
		if err != nil {
			return err
		}
		resp, err = s.executeRequest(u)
		if err != nil {
			return err
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	u, err = generateURL(s.Addr, models.Gauge, "RandomValue", strconv.FormatInt(metricSet.RandomValue, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	metricSet.Uints["PollCount"] = 0

	return nil

}

func generateURL(host, metricType, metricName, metricVal string) (*url.URL, error) {
	u, err := url.Parse("http://" + host + "/update/" + metricType + "/" + metricName + "/" + metricVal)
	if err != nil {
		return nil, err
	}
	return u, nil
}
