package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/isBlueTip/metrics/internal/agent"
	"github.com/isBlueTip/metrics/internal/agent/sender"
	"github.com/sethgrid/pester"
)

const (
	DefaultPollSeconds   = 2
	DefaultReportSeconds = 10
	DefaultRateLimit     = 3
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	serverAddr := ServerAddress{Host: "localhost", Port: "8080"}

	var pollSeconds uint
	var reportSeconds uint
	//key := ""
	var key string
	var rateLimit int

	flag.Var(&serverAddr, "a", "Server address")
	flag.UintVar(&pollSeconds, "p", DefaultPollSeconds, "Polling interval in seconds")
	flag.UintVar(&reportSeconds, "r", DefaultReportSeconds, "Reporting interval in seconds")
	flag.StringVar(&key, "k", "", "Key for SHA256 hash")
	flag.IntVar(&rateLimit, "l", DefaultRateLimit, "Num of workers to collect metrics")
	//rateLimit := flag.Int("l", DefaultRateLimit, "Num of workers to collect metrics")

	flag.Parse()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Address != nil {
		if err := serverAddr.Set(*cfg.Address); err != nil {
			panic(err)
		}
	}
	if cfg.PollInterval != nil {
		pollSeconds = *cfg.PollInterval
	}
	if cfg.ReportInterval != nil {
		reportSeconds = *cfg.ReportInterval
	}
	if cfg.Key != nil {
		key = *cfg.Key
	}
	if cfg.Rate != nil {
		rateLimit = *cfg.Rate
	}

	pollInterval := time.Duration(pollSeconds) * time.Second
	reportInterval := time.Duration(reportSeconds) * time.Second

	if err := run(serverAddr, pollInterval, reportInterval, key); err != nil {
		log.Println(err.Error())
	}
}

func run(serverAddr ServerAddress, pollInterval time.Duration, reportInterval time.Duration, key string) error {
	metric := &agent.MetricSet{
		Uints:  make(map[string]uint64),
		Floats: make(map[string]float64),
	}

	s := &sender.Sender{
		HC:   newClient(),
		Addr: net.JoinHostPort(serverAddr.Host, serverAddr.Port),
		Key:  key,
	}

	// NOTE to reviewer: Since Go 1.23, time.Tick is safe to use and garbage collected.
	pollTicker := time.Tick(pollInterval)
	reportTicker := time.Tick(reportInterval)

	for {
		select {
		case <-pollTicker:
			metric.Collect()
			log.Printf("metricSet collected: %+v\n", *metric)
		case <-reportTicker:
			//err := s.SendJSON(metric)
			err := s.SendBatch(metric)
			if err != nil {
				log.Printf("non-critical error: %s\n", err)
			}
			log.Println("metric sent")
		}
	}
}

func newClient() *pester.Client {
	client := pester.New()
	client.Timeout = 20 * time.Second
	client.MaxRetries = 3
	client.Backoff = func(retry int) time.Duration {
		intervals := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}
		if retry < len(intervals) {
			return intervals[retry]
		}
		return intervals[len(intervals)-1]
	}

	return client
}
