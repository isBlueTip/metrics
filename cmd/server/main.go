package main

import (
	"log"
	"net/http"

	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/server"
)

func run() error {
	mux := server.Router(repository.NewStorage())

	return http.ListenAndServe("127.0.0.1:8080", mux)

}

func main() {
	log.SetFlags(log.Llongfile)

	if err := run(); err != nil {
		panic(err)
	}
}
