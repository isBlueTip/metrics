package main

import (
	"log"
	"net/http"
	"time"

	"github.com/isBlueTip/metrics/internal/agent"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 4 * time.Second

const serverAddr string = "127.0.0.1:8080"

func run() error {
	metric := &agent.MetricSet{}
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
			log.Printf("metricSet collected: %+v\n", *metric)
		case <-reportTicker:
			err := sender.Send(metric)
			if err != nil {
				log.Printf("non-critical error: %s\n", err)
			}
			log.Println("metric sent")
		}
	}
}

func main() {
	log.SetFlags(log.Llongfile)

	if err := run(); err != nil {
		log.Println(err.Error())
		//panic(err)
	}
}
