package main

import (
	"errors"
	"log"
	"strconv"
)

func ValidateURL(metric []string) error {
	log.Printf("%v, type: %T, len: %v\r\n", metric, metric, len(metric))
	if len(metric) < 4 {
		return errors.New("not enough URL params to write metric, expected 4")
	}
	if metric[0] != "update" {
		return errors.New("unknown verb from URL, expected 'update'")
	}
	if metric[1] != "gauge" && metric[1] != "counter" {
		return errors.New("unknown metric type, expected 'gauge' or 'counter'")
	}
	if _, err := strconv.ParseFloat(metric[3], 64); err != nil {
		return errors.New("invalid metric value, expected numeric")
	}
	return nil
}
