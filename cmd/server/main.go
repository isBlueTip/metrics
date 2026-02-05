package main

import (
	"log"
	"net/http"

	"github.com/xxx/metrics/internal/repository"
	"github.com/xxx/metrics/internal/server"
)

func run() error {
	mux := server.Router(repository.NewStorage())

	return http.ListenAndServe("0.0.0.0:8080", mux)

}

func main() {
	log.SetFlags(log.Llongfile)

	if err := run(); err != nil {
		panic(err)
	}
}
