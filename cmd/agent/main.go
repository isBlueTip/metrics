package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/isBlueTip/metrics/internal/agent"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

func run() error {
	metric := &agent.MetricSet{PollCount: 0}
	dialer := net.Dialer{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp4", addr)
	}
	sender := &agent.Sender{
		HC: http.Client{
			Transport: transport,
		},
	}
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
