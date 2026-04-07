package main

import (
	"flag"
	"log"
	"net"
	"net/http"

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

	if err := run(&addr); err != nil {
		panic(err)
	}
}
