package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/isBlueTip/metrics/internal/logger"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/server"
	"go.uber.org/zap"
)

func run(address *ServerAddress) error {
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}

	mux := server.Router(repository.NewStorage())

	addrString := net.JoinHostPort(address.Host, address.Port)

	logger.Log.Info("Running server", zap.String("address", address.String()))

	return http.ListenAndServe(addrString, mux)

}

func main() {
	//log.SetFlags(log.Llongfile)

	addr := ServerAddress{Host: "", Port: "8080"}
	flag.Var(&addr, "a", "Parsing net address")
	flag.Parse()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Address != nil {
		if err := addr.Set(*cfg.Address); err != nil {
			panic(err)
		}
	}

	if err := run(&addr); err != nil {
		panic(err)
	}
}
