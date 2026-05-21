package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/isBlueTip/metrics/internal/agent"
	"github.com/isBlueTip/metrics/internal/hash"
	"github.com/isBlueTip/metrics/internal/models"
)

func (s *Sender) SendBatch(metricSet *agent.MetricSet) error {
	var err error
	metrics := make([]models.Update, 0)
	var metric models.Update

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
		metrics = append(metrics, metric)
	}

	for name, val := range metricSet.Floats {
		metric = models.Update{
			ID:    name,
			MType: models.Gauge,
			Value: &val,
		}
		metrics = append(metrics, metric)
	}

	ValTmp := float64(metricSet.RandomValue)
	metric = models.Update{
		ID:    "RandomValue",
		MType: models.Gauge,
		Value: &ValTmp,
	}
	metrics = append(metrics, metric)

	if len(metrics) == 0 {
		return nil
	}

	err = s.executeRequestBatch(metrics)
	if err != nil {
		return err
	}

	metricSet.Uints["PollCount"] = 0

	return nil

}

func (s *Sender) executeRequestBatch(metrics []models.Update) error {
	urlBatch := "http://" + s.Addr + "/updates/"

	log.Printf("sending metrics: %+v\n", metrics)

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	enc := json.NewEncoder(w)
	err := enc.Encode(&metrics)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlBatch, bytes.NewReader(buf.Bytes()))
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
