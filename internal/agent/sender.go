package agent

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/xxx/metrics/internal/models"
)

const serverAddr string = "127.0.0.1:8080"

type Sender struct {
	HC http.Client
}

func (s *Sender) Send(metricSet *MetricSet) error {
	var u *url.URL = nil
	var resp *http.Response = nil
	var err error = nil

	u = generateURL(models.Gauge, "Alloc", strconv.FormatUint(metricSet.Alloc, 10))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "BuckHashSys", strconv.FormatUint(metricSet.BuckHashSys, 10))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "Frees", strconv.FormatUint(metricSet.Frees, 10))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "GCCPUFraction", strconv.FormatFloat(metricSet.GCCPUFraction, 'f', 4, 64))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "GCSys", strconv.FormatFloat(metricSet.GCSys, 'f', 4, 64))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "HeapAlloc", strconv.FormatUint(metricSet.HeapAlloc, 10))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "HeapIdle", strconv.FormatUint(metricSet.HeapIdle, 10))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "HeapInuse", strconv.FormatUint(metricSet.HeapInuse, 10))
	resp, err = executeRequest(u, s.HC)
	if err != nil {
		log.Printf("error: %s\n", err)
		return err
	} else {
		log.Printf("status: %s\n", resp.Status)
		resp.Body.Close()
	}

	u = generateURL(models.Gauge, "HeapObjects", strconv.FormatUint(metricSet.HeapObjects, 10))
	resp, err = executeRequest(u, s.HC)
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

func generateURL(metricType, metricName, metricVal string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   serverAddr,
		Path:   "/update/" + metricType + "/" + metricName + "/" + metricVal + "/",
	}
}

func executeRequest(u *url.URL, hc http.Client) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	log.Printf("sending metric to url: %s\n", u.String())
	if err != nil {
		return nil, err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
