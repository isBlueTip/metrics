package main

import (
	"log"
	"net/http"
	"time"

	"github.com/isBlueTip/metrics/internal/agent"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

const serverAddr string = "127.0.0.1:8080"

func run() error {
	metric := &agent.MetricSet{PollCount: 0}
	sender := &agent.Sender{
		HC:   http.Client{},
		Addr: serverAddr,
	}

	pollTicker := time.Tick(pollInterval)
	reportTicker := time.Tick(reportInterval)

	for {
		select {
		case <-pollTicker:
			metric.Collect()
		case <-reportTicker:
			err := sender.Send(metric)
			if err != nil {
				return err
			}
		}
	}
}

func main() {
	log.SetFlags(log.Llongfile)

	if err := run(); err != nil {
		log.Fatal(err)
		//panic(err)
	}
}
