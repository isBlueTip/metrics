package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
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

func (s *Sender) Send(metricSet *MetricSet) error {
	var u *url.URL
	var err error

	for k, v := range metricSet.Uints {
		switch k {
		case "PollCount":
			u, err = generateURLWithParams(s.Addr, models.Counter, k, strconv.FormatUint(v, 10))
		default:
			u, err = generateURLWithParams(s.Addr, models.Gauge, k, strconv.FormatUint(v, 10))
		}
		if err != nil {
			return err
		}
		err = s.executeRequest(u)
		if err != nil {
			return err
		}
	}
	for k, v := range metricSet.Floats {
		u, err = generateURLWithParams(s.Addr, models.Gauge, k, strconv.FormatFloat(v, 'f', 4, 64))
		if err != nil {
			return err
		}
		err = s.executeRequest(u)
		if err != nil {
			return err
		}
	}

	u, err = generateURLWithParams(s.Addr, models.Gauge, "RandomValue", strconv.FormatInt(metricSet.RandomValue, 10))
	if err != nil {
		return err
	}
	err = s.executeRequest(u)
	if err != nil {
		return err
	}

	metricSet.Uints["PollCount"] = 0

	return nil

}

func generateURLWithParams(host, metricType, metricName, metricVal string) (*url.URL, error) {
	u, err := url.Parse("http://" + host + "/update/" + metricType + "/" + metricName + "/" + metricVal)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Sender) executeRequest(u *url.URL) error {
	log.Printf("sending metric to url: %s\n", u.String())
	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.HC.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("status: %s\n", resp.Status)

	return err
}

func (s *Sender) SendJSON(metricSet *MetricSet) error {
	var u *url.URL
	var err error
	var metric models.Update

	u, err = generateURL(s.Addr)
	if err != nil {
		return err
	}

	for name, val := range metricSet.Uints {
		switch name {
		case "PollCount":
			ValTmp := int64(val)
			metric = models.Update{
				ID:    name,
				MType: models.Counter,
				Delta: &ValTmp,
			}
		default:
			ValTmp := float64(val)
			metric = models.Update{
				ID:    name,
				MType: models.Gauge,
				Value: &ValTmp,
			}
		}
		err = s.executeRequestJSON(u, metric)
		if err != nil {
			return err
		}
	}

	for name, val := range metricSet.Floats {
		metric = models.Update{
			ID:    name,
			MType: models.Gauge,
			Value: &val,
		}
		err = s.executeRequestJSON(u, metric)
		if err != nil {
			return err
		}
	}

	u, err = generateURL(s.Addr)
	if err != nil {
		return err
	}

	ValTmp := float64(metricSet.RandomValue)
	metric = models.Update{
		ID:    "RandomValue",
		MType: models.Gauge,
		Value: &ValTmp,
	}
	err = s.executeRequestJSON(u, metric)
	if err != nil {
		return err
	}

	metricSet.Uints["PollCount"] = 0

	return nil

}

func generateURL(host string) (*url.URL, error) {
	u, err := url.Parse("http://" + host + "/update/")
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Sender) executeRequestJSON(u *url.URL, metric models.Update) error {
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	log.Printf("sending metric: %s\n", body)

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := s.HC.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//var buf bytes.Buffer
	//_, err = buf.ReadFrom(req.Body)
	//if err != nil {
	//	return err
	//}
	//log.Printf("resp: %s\n", buf)

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("status: %s\n", resp.Status)

	return err
}
