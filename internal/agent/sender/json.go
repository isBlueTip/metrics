package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/isBlueTip/metrics/internal/agent"
	"github.com/isBlueTip/metrics/internal/hash"
	"github.com/isBlueTip/metrics/internal/models"
)

func (s *Sender) SendJSON(metricSet *agent.MetricSet) error {
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
	if s.Key != "" {
		h := hash.Compute(buf.Bytes(), s.Key)
		req.Header.Set("HashSHA256", h)
	}

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
