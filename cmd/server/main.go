package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/server"
)

func run(address *ServerAddress) error {
	mux := server.Router(repository.NewStorage())

	addrString := net.JoinHostPort(address.Host, address.Port)
	return http.ListenAndServe(addrString, mux)

}

func main() {
	log.SetFlags(log.Llongfile)

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
