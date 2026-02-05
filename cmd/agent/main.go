package main

import (
	"log"
	"net/http"
	"time"

	"github.com/xxx/metrics/internal/agent"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

func run() error {
	metric := &agent.MetricSet{PollCount: 0}
	sender := &agent.Sender{HC: http.Client{}}
	start := time.Now()

	for {
		elapsed := time.Since(start).Round(time.Second)

		if elapsed%pollInterval == 0 {
			metric.Collect()
		}

		if elapsed%reportInterval == 0 {
			err := sender.Send(metric)
			if err != nil {
				return err
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	log.SetFlags(log.Llongfile)

	if err := run(); err != nil {
		panic(err)
	}
}
