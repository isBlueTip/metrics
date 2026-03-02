package agent

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/isBlueTip/metrics/internal/models"
)

type Sender struct {
	HC   http.Client
	Addr string
}

// func executeRequest(u *url.URL, s Sender) (*http.Response, error) {
func (s *Sender) executeRequest(u *url.URL) (*http.Response, error) {
	log.Printf("sending metric to url: %s\n", u.String())
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}
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

	u, err = generateURL(s.Addr, models.Gauge, "Alloc", strconv.FormatUint(metricSet.Alloc, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "BuckHashSys", strconv.FormatUint(metricSet.BuckHashSys, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "Frees", strconv.FormatUint(metricSet.Frees, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "GCCPUFraction", strconv.FormatFloat(metricSet.GCCPUFraction, 'f', 4, 64))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "GCSys", strconv.FormatUint(metricSet.GCSys, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "HeapAlloc", strconv.FormatUint(metricSet.HeapAlloc, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "HeapIdle", strconv.FormatUint(metricSet.HeapIdle, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "HeapInuse", strconv.FormatUint(metricSet.HeapInuse, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close()

	u, err = generateURL(s.Addr, models.Gauge, "HeapObjects", strconv.FormatUint(metricSet.HeapObjects, 10))
	if err != nil {
		return err
	}
	resp, err = s.executeRequest(u)
	if err != nil {
		return err
	}
	resp.Body.Close() // todo move into executeRequest

	return nil

}

func generateURL(host, metricType, metricName, metricVal string) (*url.URL, error) {
	u, err := url.Parse("http://" + host + "/update/" + metricType + "/" + metricName + "/" + metricVal)
	if err != nil {
		return nil, err
	}
	return u, nil
}
