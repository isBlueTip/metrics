package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/isBlueTip/metrics/internal/agent"
	"github.com/sethgrid/pester"
)

func run(serverAddr ServerAddress, pollInterval time.Duration, reportInterval time.Duration) error {
	metric := &agent.MetricSet{
		Uints:  make(map[string]uint64),
		Floats: make(map[string]float64),
	}

	client := pester.New()
	client.Timeout = 20 * time.Second

	sender := &agent.Sender{
		HC:   client,
		Addr: net.JoinHostPort(serverAddr.Host, serverAddr.Port),
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

	serverAddr := ServerAddress{Host: "localhost", Port: "8080"}

	var pollSeconds uint
	var reportSeconds uint

	flag.Var(&serverAddr, "a", "Server address")
	flag.UintVar(&pollSeconds, "p", 2, "Polling interval in seconds")
	flag.UintVar(&reportSeconds, "r", 10, "Reporting interval in seconds")

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
	//fmt.Printf("PollInterval: %d, type: %T\n", cfg.PollInterval, cfg.PollInterval)
	if cfg.PollInterval != nil {
		pollSeconds = *cfg.PollInterval
	}
	if cfg.ReportInterval != nil {
		reportSeconds = *cfg.ReportInterval
	}

	pollInterval := time.Duration(pollSeconds) * time.Second
	reportInterval := time.Duration(reportSeconds) * time.Second

	if err := run(serverAddr, pollInterval, reportInterval); err != nil {
		log.Println(err.Error())
		//panic(err)
	}
}
