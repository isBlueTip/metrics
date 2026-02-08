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

func (s *Sender) Send(metricSet *MetricSet) error {
	var u *url.URL = nil
	var resp *http.Response = nil
	var err error = nil

	u = generateURL(s.Addr, models.Gauge, "Alloc", strconv.FormatUint(metricSet.Alloc, 10))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "BuckHashSys", strconv.FormatUint(metricSet.BuckHashSys, 10))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "Frees", strconv.FormatUint(metricSet.Frees, 10))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "GCCPUFraction", strconv.FormatFloat(metricSet.GCCPUFraction, 'f', 4, 64))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "GCSys", strconv.FormatFloat(metricSet.GCSys, 'f', 4, 64))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "HeapAlloc", strconv.FormatUint(metricSet.HeapAlloc, 10))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "HeapIdle", strconv.FormatUint(metricSet.HeapIdle, 10))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "HeapInuse", strconv.FormatUint(metricSet.HeapInuse, 10))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(s.Addr, models.Gauge, "HeapObjects", strconv.FormatUint(metricSet.HeapObjects, 10))
	resp, err = executeRequest(u, *s)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	log.Printf("metric sent\n")
	return nil

}

func generateURL(host, metricType, metricName, metricVal string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/update/" + metricType + "/" + metricName + "/" + metricVal,
	}
}

func executeRequest(u *url.URL, s Sender) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	log.Printf("sending metric to url: %s\n", u.String())
	if err != nil {
		return nil, err
	}
	resp, err := s.HC.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
